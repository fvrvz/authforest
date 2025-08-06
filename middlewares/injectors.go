package middlewares

import (
	"auth-service-go/initializers"

	"github.com/gin-gonic/gin"
)

func InjectDeps(cfg *initializers.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("config", cfg)
		ctx.Next()
	}
}
