package controllers

import (
	"github.com/fvrvz/auth-service-go/services"
	"github.com/gin-gonic/gin"
)

func SetupOIDCRoutes(router *gin.Engine) {
	// OIDC Discovery (must be at the root, not under /api/v1)
	router.GET("/.well-known/openid-configuration", services.OIDCDiscovery)
	router.GET("/.well-known/jwks.json", services.JWKS)

	// OAuth2 endpoints
	oauth2 := router.Group("/oauth2")
	{
		oauth2.GET("/authorize", services.Authorize)
		oauth2.POST("/authorize", services.HandleLogin)
		oauth2.POST("/token", services.TokenExchange)
		oauth2.GET("/userinfo", services.UserInfo)
		oauth2.POST("/userinfo", services.UserInfo)
	}

	// Client registration (protected - behind auth middleware in router.go)
}

func SetupOIDCPrivateRoutes(routerGroup *gin.RouterGroup) {
	oauth2 := routerGroup.Group("/oauth2")
	{
		oauth2.POST("/register", services.RegisterClient)
		oauth2.GET("/clients", services.ListClients)
		oauth2.GET("/clients/:clientId", services.GetClient)
		oauth2.PATCH("/clients/:clientId", services.UpdateClient)
		oauth2.DELETE("/clients/:clientId", services.DeleteClient)
		oauth2.GET("/stats", services.DashboardStats)
	}
}
