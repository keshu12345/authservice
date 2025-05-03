package service

import (
	"fmt"

	"github.com/keshucs12345/authservice/models"
	"github.com/keshucs12345/authservice/utils"
	"golang.org/x/crypto/bcrypt"
)

func (a *authService) Signup(username, password, role string) (bool, error) {
	user, _ := a.AuthDAO.GetUserByUsername(username)

	if user != nil {
		a.Logger.Errorf("Username already exists: %s", username)
		return false, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.Logger.Errorf("Failed to hash password: %v", err)
		return false, fmt.Errorf("failed to hash password: %v", err)
	}

	apiKey, err := utils.GenerateSecureAPIKey(32)
	if err != nil {
		a.Logger.Errorf("Failed to generate API key: %v", err)
		return false, fmt.Errorf("failed to generate API key: %v", err)
	}

	// Create new user
	newUser := &models.User{
		UserID:       utils.GenerateUserID(role),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		APIKey:       apiKey,
	}

	err = a.AuthDAO.InsertUser(newUser)
	if err != nil {
		a.Logger.Errorf("Failed to insert user: %v", err)
		return false, fmt.Errorf("failed to insert user: %v", err)
	}
	a.Logger.Infof("User registered successfully: %s", username)
	return true, nil
}

func (a *authService) SignIn(username string, password string) (*models.User, string, error) {
	user, err := a.AuthDAO.GetUserByUsername(username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		a.Logger.Errorf("Login attempt failed for non-existent user: %s", username)
		return nil, "", fmt.Errorf("user not found")

	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		a.Logger.Errorf("user not found: %s", username)
		return nil, "", fmt.Errorf("invalid username or password")

	}
	// Generate a new token
	token, err := utils.GenerateJWT(user, a.Cfg.JwtSecretKey,a.Cfg.Issuer, a.Cfg.TokenDuration)
	if err != nil {
		a.Logger.Errorf("Failed to generate token for user %s: %v", username, err)
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}
	a.Logger.Infof("User logged in successfully: %s", username)
	return user, token, nil
}

func (a *authService) RefreshToken(token string) (string, error) {
	// Parse the token
	claims, err := utils.ParseToken(token, a.Cfg.JwtSecretKey)
	if err != nil {
		a.Logger.Errorf("Failed to parse token: %v", err)
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Get User by Username
	user, err := a.AuthDAO.GetUserByUsername(claims.Subject)
	if err != nil {
		a.Logger.Errorf("Failed to get user by username: %v", err)
		return "", fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		a.Logger.Errorf("User not found for username: %s", claims.Subject)
		return "", fmt.Errorf("user not found")
	}

	// Generate a new token
	newToken, err := utils.GenerateJWT(user, a.Cfg.JwtSecretKey, a.Cfg.Issuer, a.Cfg.TokenDuration)
	if err != nil {
		a.Logger.Errorf("Failed to generate new token: %v", err)
		return "", fmt.Errorf("failed to generate new token: %v", err)
	}
	a.Logger.Infof("Token refreshed successfully for user: %s", claims.Subject)
	return newToken, nil
}

func (a *authService) RevokeToken(token string) error {

	calims, err := utils.ParseToken(token, a.Cfg.JwtSecretKey)
	if err != nil {
		a.Logger.Errorf("Failed to parse token: %v", err)
		return fmt.Errorf("failed to parse token: %v", err)
	}
	user, err := a.AuthDAO.GetUserByUsername(calims.Subject)
	if err != nil {
		a.Logger.Errorf("Failed to get user by username: %v", err)
		return fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		a.Logger.Errorf("User not found for username: %s", calims.Subject)
		return fmt.Errorf("user not found")
	}
	user.IsRevoke = true
	a.Logger.Infof("Token revoked successfully for user: %s", calims.Subject)
	return a.AuthDAO.UpdateUser(user)
}

func (a *authService) UnRevokeToken(token string) error {
	calims, err := utils.ParseToken(token, a.Cfg.JwtSecretKey)
	if err != nil {
		a.Logger.Errorf("Failed to parse token: %v", err)
		return fmt.Errorf("failed to parse token: %v", err)
	}
	user, err := a.AuthDAO.GetUserByUsername(calims.Subject)
	if err != nil {
		a.Logger.Errorf("Failed to get user by username: %v", err)
		return fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		a.Logger.Errorf("User not found for username: %s", calims.Subject)
		return fmt.Errorf("user not found")
	}
	user.IsRevoke = false
	a.Logger.Infof("Token unrevoked successfully for user: %s", calims.Subject)

	return a.AuthDAO.UpdateUser(user)
}
