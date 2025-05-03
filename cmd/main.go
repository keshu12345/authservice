package main

import (
	"flag"
	"log"
	"os"

	"github.com/keshucs12345/authservice/config"
	"github.com/keshucs12345/authservice/dao"
	"github.com/keshucs12345/authservice/db"
	"github.com/keshucs12345/authservice/logger"
	"github.com/keshucs12345/authservice/migrate"
	"github.com/keshucs12345/authservice/pkg"
	"github.com/keshucs12345/authservice/server/router"
	"go.uber.org/fx"
)

var configDirPath = flag.String("config", "", "path for config dir")
var migrationDir = flag.String("migrations", "./migrations", "path for migration dir")

func main() {

	flag.Parse()
	log.New(os.Stdout, "", 0)
	app := fx.New(
		config.NewFxModule(*configDirPath, ""),
		router.Module,
		pkg.Module,
		db.Module,
		migrate.Module(*migrationDir),
		dao.Module,
		logger.Module,
	)

	app.Run()
}
