package middlewares

import (
	"net/http"

	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := helpers.ExtractTokenFromHeaders(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		claims, err := helpers.Verify(token)

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
