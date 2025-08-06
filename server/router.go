package server

import (
	"github.com/fvrvz/auth-service-go/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		controllers.SetupAuthRoutes(v1)
		controllers.SetupUserRoutes(v1)
	}

	return router
}
