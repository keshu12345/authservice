package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/keshucs12345/authservice/config"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type Postgres struct {
	Db *sqlx.DB
	fx.Out
}

func NewDBInstance(conf *config.Configuration) (Postgres, error) {

	log.Info("DB connected")

	db, err := sqlx.Connect("postgres", dataSource(&conf.DB))
	if err != nil {
		log.WithField("error", err).Error("Database Error")
		return Postgres{}, err
	}

	// Check if is alive
	_, err = CheckConnection(db)
	if err != nil {
		log.WithField("error", err).Error("Database Error")
		return Postgres{}, err
	}
	db.SetMaxIdleConns(conf.DB.MaxIdle)
	db.SetMaxOpenConns(conf.DB.MaxOpen)

	return Postgres{
		Db: db,
	}, nil
}

func NewDBSQLInstance(dbConf *config.DB) (db *sql.DB, err error) {
	return sql.Open(dbConf.Driver, dataSource(dbConf))
}

// DataSource returns config in format required by SQL
// "postgres://<user>:<password>@<host>/<database>?options"
func dataSource(dbConf *config.DB) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?%s", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Database, dbConf.Options)
}

func CheckConnection(sql *sqlx.DB) (bool, error) {
	log.Info("Check connection....")
	if sql == nil {
		return false, errors.New("connection not found")
	}
	err := sql.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}
