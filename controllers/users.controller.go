package controllers

import (
	"github.com/fvrvz/auth-service-go/services"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("/", services.GetUsers)
		users.GET("/:userId")
		users.POST("/register", services.Register)
		users.DELETE("/:userId", services.Delete)
	}
}
