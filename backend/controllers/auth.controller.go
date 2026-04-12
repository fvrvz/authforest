package controllers

import (
	"github.com/fvrvz/authforest/services"
	"github.com/gin-gonic/gin"
)

func SetupAuthPublicRoutes(routerGroup *gin.RouterGroup) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/login", services.Login)
	}
}

func SetupAuthPrivateRoutes(routerGroup *gin.RouterGroup) {
	auth := routerGroup.Group("/auth")
	{
		auth.GET("/logout", services.Logout)
		auth.POST("/refresh", services.RotateRefreshToken)
	}
}
