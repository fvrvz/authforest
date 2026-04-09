package middlewares

import (
	"net/http"

	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

		// Try HS256 (legacy) first, then RS256 (OIDC access tokens)
		var claims *jwt.RegisteredClaims

		claims, err = helpers.ExtractClaims(token)
		if err != nil {
			// Fallback: try RS256 OIDC access token
			claims, err = helpers.ExtractRSAClaims(token)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "Invalid or expired token",
				})
				return
			}
		}

		var blacklistedRecord models.AccessTokenBlacklist
		if err := db.GetDB().First(&blacklistedRecord, models.AccessTokenBlacklist{JTI: claims.ID}).Error; err == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Token is expired",
			})
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
