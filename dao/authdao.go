package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/keshucs12345/authservice/logger"
	"github.com/keshucs12345/authservice/models"
	"go.uber.org/fx"
)

type AuthDAO interface {
	InsertUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
	UpdateUser(user *models.User) error
	GetUserByAPIKey(apiKey string) (*models.User, error)
}

type authDAO struct {
	DB     *sqlx.DB
	Logger logger.Logger
}

type authDAOParams struct {
	fx.In
	DB     *sqlx.DB
	Logger logger.Logger
}

func New(adp authDAOParams) AuthDAO {
	return authDAO{
		DB:     adp.DB,
		Logger: adp.Logger,
	}
}
