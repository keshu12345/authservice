package endpoint

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/keshucs12345/authservice/models"
	"github.com/keshucs12345/authservice/utils"
)

func Signup(c *gin.Context) {

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appLogger.Errorf("failed to bind request: %v", err)
		utils.RespondError(c, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if req.Username == "" || len(req.Password) < 8 || (req.Role != "creator" && req.Role != "partner") {
		appLogger.Errorf("invalid input: username required, password min 8 chars, role must be 'creator' or 'partner': request %+v", req)
		utils.RespondError(c, http.StatusBadRequest, fmt.Errorf("invalid input: username required, password min 8 chars, role must be 'creator' or 'partner'"))
		return
	}

	registered, err := authService.Signup(req.Username, req.Password, req.Role)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			appLogger.Warnf("username already exists: %s", req.Username)
			utils.RespondError(c, http.StatusConflict, fmt.Errorf("username already exists"))
		} else {
			appLogger.Errorf("failed to register user: %v", err)
			utils.RespondError(c, http.StatusInternalServerError, fmt.Errorf("failed to register user: %w", err))
		}
		return
	}

	if !registered {
		appLogger.Warnf("username already exists: %s", req.Username)
		utils.RespondError(c, http.StatusConflict, fmt.Errorf("username already exists"))
		return
	}

	registeredUser := models.User{
		Username: req.Username,
	}

	utils.RespondSuccess(c, http.StatusCreated, "User registered successfully", registeredUser)
	appLogger.Infof("User registered successfully: %s", req.Username)
}

func SignIn(c *gin.Context) {

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appLogger.Errorf("failed to bind request: %v", err)
		utils.RespondError(c, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if req.Username == "" || req.Password == "" {
		appLogger.Errorf("invalid input: username and password required: request %+v", req)
		utils.RespondError(c, http.StatusBadRequest, fmt.Errorf("username and password required"))
		return
	}

	user, token, err := authService.SignIn(req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid username or password") {
			appLogger.Warnf("login attempt failed for non-existent user: invalid username or password %s", req.Username)
			utils.RespondError(c, http.StatusUnauthorized, err)
		} else if strings.Contains(err.Error(), "user not found") {
			appLogger.Warnf("login attempt failed for non-existent user: %s", req.Username)
			utils.RespondError(c, http.StatusNotFound, err)
		} else {
			appLogger.Errorf("failed to sign in user: %v", err)
			utils.RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	login := models.LoginResponse{
		Token:  token,
		APIKey: user.APIKey,
		UserID: user.UserID,
		Role:   user.Role,
	}

	utils.RespondSuccess(c, http.StatusOK, "Login successful", login)
	appLogger.Infof("User logged in successfully: %+v", login)
}

func RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appLogger.Errorf("failed to bind request: %v", err)
		utils.RespondError(c, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if req.Token == "" {
		appLogger.Warnln("token can not be empty")
		utils.RespondError(c, http.StatusBadRequest, fmt.Errorf("token can not be empty"))
		return
	}

	newToken, err := authService.RefreshToken(req.Token)
	if err != nil {
		if strings.Contains(err.Error(), "invalid token") {
			appLogger.Warnf("token is invalid: %s", req.Token)
			utils.RespondError(c, http.StatusUnauthorized, err)
		} else {
			appLogger.Errorf("failed to refresh token: %v", err)
			utils.RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	response := models.RefreshTokenResponse{
		Message:      "Token refreshed successfully.",
		RefreshToken: newToken,
	}

	utils.RespondSuccess(c, http.StatusOK, "Token refreshed successfully", response)
	appLogger.Infof("Token refreshed successfully: %s", newToken)
}

// RevokeToken handler for Gin (revokes until actual token expiry)
func RevokeToken(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		appLogger.Warnln("missing or malformed Authorization header")
		utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf("missing or malformed Authorization header"))
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := authService.RevokeToken(token); err != nil {
		if strings.Contains(err.Error(), "token is invalid") {
			appLogger.Warnf("token is invalid: %s", token)
			utils.RespondError(c, http.StatusUnauthorized, err)
		} else {
			appLogger.Errorf("failed to revoke token: %v", err)
			utils.RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.RespondSuccess(c, http.StatusAccepted, "Token revoked successfully", nil)
	appLogger.Infof("Token revoked successfully: %s", token)
}

func UnRevokeToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		appLogger.Warnln("missing or malformed Authorization header")
		utils.RespondError(c, http.StatusUnauthorized, fmt.Errorf("missing or malformed Authorization header"))
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := authService.UnRevokeToken(token); err != nil {
		if strings.Contains(err.Error(), "token is invalid") {
			appLogger.Warnf("token is invalid: %s", token)
			utils.RespondError(c, http.StatusUnauthorized, err)
		} else {
			appLogger.Errorf("failed to un-revoke token: %v", err)
			utils.RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Token unrevoked successfully", nil)
	appLogger.Infof("Token unrevoked successfully: %s", token)
}
