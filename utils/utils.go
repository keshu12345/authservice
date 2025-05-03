package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/keshucs12345/authservice/constants"
	"github.com/keshucs12345/authservice/models"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSecureAPIKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateUserID(role string) string {
	return constants.User + role + "-" + time.Now().Format(constants.Formatter)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(user *models.User, jwtSecretKey string, issuer string, tokenDuration int) (string, error) {

	expirationTime := time.Now().Add(time.Duration(tokenDuration) * time.Minute)
	tokenID := fmt.Sprintf("%x", time.Now().UnixNano())

	claims := &models.CustomClaims{
		UserID: user.UserID,
		Roles:  []string{user.Role},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString, jwtSecretKey string) (*models.CustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid or expired token")
}

func CheckRoles(userRoles []string, requiredRoles []string) bool {
	if len(requiredRoles) == 0 {
		return true
	}
	userRoleSet := make(map[string]struct{})
	for _, role := range userRoles {
		userRoleSet[role] = struct{}{}
	}
	for _, requiredRole := range requiredRoles {
		if _, found := userRoleSet[requiredRole]; found {
			return true
		}
	}
	return false
}
