package controllers

import (
	"github.com/fvrvz/auth-service-go/services"
	"github.com/gin-gonic/gin"
)

func SetupUserPrivateRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("", services.GetUsers)
		users.GET("/:userId", services.GetUser)
		users.DELETE("/:userId", services.Delete)
	}
}

func SetupUserPublicRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("/register", services.Register)
	}
}
