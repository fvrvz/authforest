package controllers

import (
	"github.com/fvrvz/auth-service-go/services"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(routerGroup *gin.RouterGroup) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/login", services.Login)
		auth.POST("/logout", services.Logout)
		auth.POST("/refresh", services.RefreshToken)
	}
}
