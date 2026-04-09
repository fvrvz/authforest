package controllers

import (
	"github.com/fvrvz/auth-service-go/services"
	"github.com/gin-gonic/gin"
)

func SetupUserPrivateRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("", services.GetUsers)
		users.POST("/create", services.AdminCreateUser)
		users.GET("/:userId", services.GetUser)
		users.DELETE("/:userId", services.Delete)
		users.PATCH("/:userId", services.UpdateUser)
		users.PATCH("/:userId/admin", services.AdminUpdateUser)
	}

	roles := r.Group("/roles")
	{
		roles.GET("", services.ListRoles)
		roles.POST("", services.CreateRole)
		roles.GET("/:roleId", services.GetRole)
		roles.PATCH("/:roleId", services.UpdateRole)
		roles.DELETE("/:roleId", services.DeleteRole)
	}
}

func SetupUserPublicRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("/register", services.Register)
		users.POST("/request-password-reset", services.RequestPasswordReset)
		users.POST("/reset-password", services.ResetPassword)
	}
}
