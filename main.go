package main

import (
	"auth-service-go/controllers"
	"auth-service-go/initializers"
	"auth-service-go/middlewares"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := initializers.LoadConfig("config/config.yml", ".env")

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	initializers.InitDB(cfg)

	// set routes
	router := gin.Default()

	router.Use(middlewares.InjectDeps(cfg))

	v1 := router.Group("/api/v1")
	{
		controllers.SetupAuthRoutes(v1)
		controllers.SetupUserRoutes(v1)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
