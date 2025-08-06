package helpers

import (
	"auth-service-go/initializers"

	"github.com/gin-gonic/gin"
)

func GetConfig(ctx *gin.Context) *initializers.Config {
	val, exists := ctx.Get("config")

	if !exists {
		panic("Config not found in the context")
	}

	return val.(*initializers.Config)
}
