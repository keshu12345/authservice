package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keshucs12345/authservice/config"
	"github.com/keshucs12345/authservice/server"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewGinRouter,
	),
	fx.Invoke(
		server.Initialize,
	),
)

func NewGinRouter(config *config.Configuration) (*gin.Engine, error) {
	g := gin.New()
	g.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length", "Content-Type", "Last-Modified"},
			AllowCredentials: true,
			MaxAge:           1 * time.Hour,
		},
		),
	)
	return g, nil
}
