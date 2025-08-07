package server

import (
	"time"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/controllers"
	"github.com/fvrvz/auth-service-go/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	setupCors(router)

	v1 := router.Group("/api/v1")
	{
		controllers.SetupAuthRoutes(v1)
		controllers.SetupUserPublicRoutes(v1)

		privateGroup := v1.Group("/", middlewares.AuthMiddleware())
		controllers.SetupUserPrivateRoutes(privateGroup)
	}

	return router
}

func setupCors(router *gin.Engine) {
	corsConfig := config.GetConfig().Server.CORS
	corsDefaults := cors.DefaultConfig()

	if len(corsConfig.AllowOrigins) > 0 {
		corsDefaults.AllowOrigins = corsConfig.AllowOrigins
	}
	if len(corsConfig.AllowMethods) > 0 {
		corsDefaults.AllowMethods = corsConfig.AllowMethods
	}
	if len(corsConfig.AllowHeaders) > 0 {
		corsDefaults.AllowHeaders = corsConfig.AllowHeaders
	}
	if len(corsConfig.ExposeHeaders) > 0 {
		corsDefaults.ExposeHeaders = corsConfig.ExposeHeaders
	}
	if corsConfig.MaxAge > 0 {
		corsDefaults.MaxAge = time.Duration(corsConfig.MaxAge) * time.Hour
	}

	router.Use(cors.New(corsDefaults))
}
