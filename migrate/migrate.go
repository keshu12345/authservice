package migrate

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func Module(migrationDir string) fx.Option {
	return fx.Invoke(func(db *sqlx.DB) {
		log.Debug("Setting Goose dialect...")
		if err := goose.SetDialect("postgres"); err != nil {
			log.WithError(err).Fatal("Failed to set dialect")
		}

		log.Info("Running DB migrations...")
		if err := goose.Up(db.DB, migrationDir); err != nil {
			log.WithError(err).Fatal("Migration failed")
		}
		log.Info("Migrations applied s	uccessfully")
	})
}
