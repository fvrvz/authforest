package controllers

import (
	"auth-service-go/services"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(routerGroup *gin.RouterGroup) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/register", services.Register)
		auth.POST("/login", services.Login)
		auth.DELETE("/delete", services.Delete)
	}
}
