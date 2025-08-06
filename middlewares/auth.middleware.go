package middlewares

import (
	"net/http"

	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := services.Verify(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
