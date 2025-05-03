package service

import (
	"fmt"

	"github.com/keshucs12345/authservice/models"
	"github.com/keshucs12345/authservice/utils"
)

func (a *authService) ApiKeyAuthenticate(apiKey string) (*models.User, error) {
	user, err := a.AuthDAO.GetUserByAPIKey(apiKey)
	if err != nil {
		a.Logger.Errorf("failed to get user by API key: %v", err)
		return nil, fmt.Errorf("failed to get user by API key: %v", err)
	}
	if user == nil {
		a.Logger.Warnf("User does not exit: %s", apiKey)
		return nil, fmt.Errorf("user does not exit: %s", apiKey)
	}

	a.Logger.Infof("User authenticated successfully with API Key :%s with username : %s", apiKey, user.Username)
	return user, nil
}

func (a *authService) JwtTokenAuthenticate(token string) (*models.CustomClaims, bool, error) {
	claims, err := utils.ParseToken(token, a.Cfg.JwtSecretKey)

	a.Logger.Infof(":::::::::::::Claims::::::::::::: %+v", claims)
	if err != nil {
		a.Logger.Errorf("failed to parse token: %v", err)
		return nil, false, err
	}

	// if claims.ExpiresAt.Time.Before(time.Now()) {
	// 	a.Logger.Warnf("token expired: %v", claims.ExpiresAt.Time)
	// 	return nil, false, fmt.Errorf("token expired")
	// }

	user, err := a.AuthDAO.GetUserByID(claims.UserID)
	a.Logger.Infof(":::::::::::::User::::::::::::: %+v", user)
	if err != nil {
		a.Logger.Errorf("failed to get user by ID: %v", err)
		return nil, false, fmt.Errorf("failed to get user by ID: %v", err)
	}
	if user == nil {
		a.Logger.Warnf("user not found: %s", claims.UserID)
		return nil, false, fmt.Errorf("user not found")
	}

	a.Logger.Infof("User authenticated successfully with token: %s with username: %s", token, user.Username)
	return claims, user.IsRevoke, nil
}
