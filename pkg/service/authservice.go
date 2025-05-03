package service

import (
	"github.com/keshucs12345/authservice/config"
	"github.com/keshucs12345/authservice/dao"
	"github.com/keshucs12345/authservice/logger"
	"github.com/keshucs12345/authservice/models"
	"go.uber.org/fx"
)

type Authenticator interface {
	Signup(username, password, role string) (bool, error)
	SignIn(username string, password string) (*models.User, string, error)
	RefreshToken(token string) (string, error)
	RevokeToken(token string) error
	UnRevokeToken(token string) error
}

type MiddlewareAuthenticator interface {
	ApiKeyAuthenticate(apiKey string) (*models.User, error)
	JwtTokenAuthenticate(token string) (*models.CustomClaims, bool, error)
}

type AuthService interface {
	Authenticator
	MiddlewareAuthenticator
}

type authService struct {
	AuthDAO dao.AuthDAO
	Cfg     *config.Configuration
	Logger  logger.Logger
}

type authServiceParams struct {
	fx.In
	AuthDAO dao.AuthDAO
	Config  *config.Configuration
	Logger  logger.Logger
}

func NewAuthService(asp authServiceParams) AuthService {
	return &authService{
		AuthDAO: asp.AuthDAO,
		Cfg:     asp.Config,
		Logger:  asp.Logger,
	}
}
