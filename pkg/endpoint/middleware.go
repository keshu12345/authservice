package endpoint

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/keshucs12345/authservice/constants"
	"github.com/keshucs12345/authservice/models"
	"github.com/keshucs12345/authservice/utils"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		clientIP := c.ClientIP()

		appLogger.Infof("--> [%s] %s %s %s", clientIP, c.Request.Method, c.Request.URL.Path, c.Request.Proto)

		c.Next()

		clientInfo := ""
		if clientName, exists := c.Get(constants.ApiKeyUserCtxKey.String()); exists {
			if nameStr, ok := clientName.(string); ok {
				clientInfo = fmt.Sprintf(" (Client: %s)", nameStr)
			}
		}
		appLogger.Infof("clientInfo: %s", clientInfo)

		userInfo := ""
		if claims, exists := c.Get(constants.ApiKeyClientCtxKey.String()); exists {
			if cc, ok := claims.(*models.CustomClaims); ok && cc != nil {
				userInfo = fmt.Sprintf(" (UserID: %s, Roles: %v)", cc.UserID, cc.Roles)
			}
		}
		appLogger.Infof("<-- [%s] %s %s %s%s%s (%s)", clientIP, c.Request.Method, c.Request.URL.Path, c.Request.Proto, clientInfo, userInfo, time.Since(start))
	}
}

func ApiKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			appLogger.Warnf("API Key not found in header")
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf(" [ApiKeyAuthMiddleware] missing API Key"))
			c.Abort()
			return
		}

		user, err := authService.ApiKeyAuthenticate(apiKey)
		if err != nil {
			appLogger.Warnf("API Key authentication failed: %v", err)
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf(" [ApiKeyAuthMiddleware] invalid API Key: %w", err))
			c.Abort()
			return
		}

		c.Set(constants.ApiKeyUserCtxKey.String(), user)
		c.Next()
		appLogger.Infof("API Key authentication successful for user: %s", user.Username)
	}
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appLogger.Infof("JWT Auth Error: Missing Authorization header for %s %s", c.Request.Method, c.Request.URL.Path)
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf("[JwtAuthMiddleware] missing Authorization header"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			appLogger.Infof("JWT Auth Error: Invalid Authorization format for %s %s", c.Request.Method, c.Request.URL.Path)
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf(" [JwtAuthMiddleware] invalid Authorization format"))
			c.Abort()
			return
		}

		token := parts[1]

		claims, isRevoked, err := authService.JwtTokenAuthenticate(token)
		if err != nil {
			appLogger.Infof("JWT Auth Error: Token validation failed for %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
			errMsg := " [JwtAuthMiddleware] Unauthorized: Invalid token"
			if errors.Is(err, jwt.ErrTokenExpired) {
				errMsg = "[JwtAuthMiddleware] Unauthorized: Token expired"
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				errMsg = " [JwtAuthMiddleware] Unauthorized: Invalid token signature"
			}
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf("%w", errors.New(errMsg)))
			c.Abort()
			return
		}

		serviceRole := path.Base(c.FullPath())
		appLogger.Infof("Service Role: %s", serviceRole)

		requiredRoles, exist := constants.RouteRoles["creators"]
		if !exist {
			appLogger.Infof("JWT Auth Error: User %s not found in registered roles", claims.UserID)
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf(" [JwtAuthMiddleware] unauthorized: User not found"))
			c.Abort()
			return
		}
		if len(requiredRoles) == 0 {
			appLogger.Infof("JWT Auth Error: No roles found for user %s", claims.UserID)
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf(" [JwtAuthMiddleware] unauthorized: No roles found for user"))
			c.Abort()
			return
		}

		if !utils.CheckRoles(claims.Roles, requiredRoles) {
			clientUsername := "N/A"
			appLogger.Infof("JWT Auth Error: User %s does not have required roles %v for %s %s", claims.Subject, requiredRoles, c.Request.Method, c.Request.URL.Path)
			if user, ok := c.Request.Context().Value(constants.ApiKeyClientCtxKey.String()).(models.User); ok {
				clientUsername = user.Username
			}
			appLogger.Infof("RBAC Error: User '%s' via JWT (Client: %s) with roles %v does not have required roles %v for %s %s",
				claims.Subject, clientUsername, claims.Roles, requiredRoles, c.Request.Method, c.Request.URL.Path)

			utils.RespondError(c, http.StatusForbidden, fmt.Errorf(" [JwtAuthMiddleware] forbidden: Insufficient permissions"))
			c.Abort()
			return
		}

		if isRevoked {
			appLogger.Infof("JWT Auth Error: Token revoked for %s %s", c.Request.Method, c.Request.URL.Path)
			utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf(" [JwtAuthMiddleware] unauthorized: Token revoked"))
			c.Abort()
			return
		}

		c.Set(constants.CaimsContextKey.String(), claims)
		c.Next()
		appLogger.Infof("JWT Auth: User %s with roles %v authenticated successfully for %s %s", claims.Subject, claims.Roles, c.Request.Method, c.Request.URL.Path)
	}
}
