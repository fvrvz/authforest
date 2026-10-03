package controllers

import (
	"net/http"

	"github.com/fvrvz/authforest/db"
	"github.com/gin-gonic/gin"
)

// SetupHealthRoutes registers the health check endpoint.
// The /health route is intentionally outside the /api/v1 prefix and
// requires no authentication so orchestrators can probe it freely.
func SetupHealthRoutes(router *gin.Engine) {
	router.GET("/health", healthCheck)
}

func healthCheck(c *gin.Context) {
	sqlDB, err := db.GetDB().DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "failed to get database handle",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database unreachable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
