package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/keshucs12345/authservice/logger"
	"github.com/keshucs12345/authservice/pkg/service"
)

var authService service.AuthService

var appLogger logger.Logger

func RegisterEndpoint(g *gin.Engine, as service.AuthService, logger logger.Logger) {
	appLogger = logger
	authService = as

	auth := g.Group("/v1")
	{
		auth.Use(LoggingMiddleware())

		auth.POST("/signin", SignIn)
		auth.POST("/signup", Signup)
		auth.POST("/refresh", RefreshToken)

		auth.POST("/revoke", RevokeToken)
		auth.POST("/unrevoke", UnRevokeToken)
	}

	api := g.Group("/api")
	{
		api.Use(LoggingMiddleware())
		api.Use(ApiKeyAuthMiddleware())
		api.Use(JwtAuthMiddleware())

		api.GET("/creators", Creator)
	}
}
