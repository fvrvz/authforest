package server

import (
	"github.com/fvrvz/auth-service-go/controllers"
	"github.com/fvrvz/auth-service-go/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		controllers.SetupAuthRoutes(v1)
		controllers.SetupUserPublicRoutes(v1)

		privateRoutes := v1.Group("/", middlewares.AuthMiddleware())
		controllers.SetupUserPrivateRoutes(privateRoutes)
	}

	return router
}
