package services

import (
	"log"
	"net/http"
	"time"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/constants"
	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Invalid Input",
			Description: err.Error(),
		})
		return
	}

	user, err := authenticateUser(req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := generateTokensAndStoreRefreshToken(user.Username)

	if err != nil {
		log.Printf("Login token error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Token generation failed"})
		return
	}

	respondWithTokens(c, accessToken, refreshToken)
}

func Logout(ctx *gin.Context) {
	accessToken, err := helpers.ExtractTokenFromHeaders(ctx)

	if err != nil {
		log.Fatalf("token is missing %v", err.Error())
		return
	}

	claims, err := helpers.ExtractClaims(accessToken)

	if err != nil {
		log.Fatalf("token expired or invalid >> %v", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Invalid or expired refresh token",
			Description: err.Error(),
		})
		return
	}

	blacklistRecord := models.AccessTokenBlacklist{
		JTI:       claims.ID,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}

	if err := db.GetDB().FirstOrCreate(&blacklistRecord, models.AccessTokenBlacklist{JTI: claims.ID}).Error; err != nil {
		log.Printf("Failed to blacklist the access token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Failed to logout",
			Description: err.Error(),
		})
		return
	}

	//delete the refresh token associated with this accessToken claim ID
	if err := db.GetDB().Delete(models.AuthRefreshTokens{}, models.AuthRefreshTokens{AccessTokenID: claims.ID}).Error; err != nil {
		log.Printf("Failed to delete the refresh token associated with access token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Failed to logout",
			Description: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[string]{
		Message: "User logged out successfully",
	})
}

func RotateRefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Invalid Input",
			Description: err.Error(),
		})
		return
	}

	refreshTokenClaims, err := helpers.ExtractClaims(req.RefreshToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Invalid or expired refresh token",
			Description: err.Error(),
		})
		return
	}

	if err := db.GetDB().Where("jti = ? AND expires_at > ?", refreshTokenClaims.ID, time.Now()).Delete(models.AuthRefreshTokens{}).Error; err != nil {
		log.Printf("Failed to delete refresh token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
		return
	}

	accessToken, _ := helpers.ExtractTokenFromHeaders(ctx)

	accessTokenClaims, err := helpers.ExtractClaims(accessToken)

	if err == nil {
		// black list the acccess token if not expired
		blacklistRecord := models.AccessTokenBlacklist{
			JTI:       accessTokenClaims.ID,
			ExpiresAt: accessTokenClaims.ExpiresAt.Time,
			IssuedAt:  accessTokenClaims.IssuedAt.Time,
		}

		if err := db.GetDB().FirstOrCreate(&blacklistRecord, models.AccessTokenBlacklist{JTI: accessTokenClaims.ID}).Error; err != nil {
			log.Printf("Failed to blacklist the access token: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:       "Failed to logout",
				Description: err.Error(),
			})
			return
		}
	}

	newAccessToken, newRefreshToken, err := generateTokensAndStoreRefreshToken(refreshTokenClaims.Subject)

	if err != nil {
		log.Printf("Failed to generate new tokens: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Token generation failed",
			Description: err.Error(),
		})
		return
	}

	respondWithTokens(ctx, newAccessToken, newRefreshToken)
}

func generateTokensAndStoreRefreshToken(username string) (string, string, error) {
	config := config.GetConfig()

	accessToken, refreshToken, err := helpers.GenerateJWT(config, username)

	if err != nil {
		return "", "", err
	}

	accessTokenClaims, err := helpers.ExtractClaims(accessToken)

	if err != nil {
		return "", "", err
	}

	refreshTokenClaims, err := helpers.ExtractClaims(refreshToken)

	if err != nil {
		return "", "", err
	}

	newRecord := models.AuthRefreshTokens{
		JTI:           refreshTokenClaims.ID,
		Username:      username,
		IssuedAt:      refreshTokenClaims.IssuedAt.Time,
		ExpiresAt:     refreshTokenClaims.ExpiresAt.Time,
		AccessTokenID: accessTokenClaims.ID,
	}

	if err := db.GetDB().Create(&newRecord).Error; err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func respondWithTokens(ctx *gin.Context, accessToken string, refreshToken string) {
	ctx.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    constants.BEARER,
		ExpiresIn:    config.GetConfig().JWT.ExpiryMinutes * 60,
	})
}

func authenticateUser(req dto.LoginRequest) (*models.User, error) {
	var user models.User

	if err := db.GetDB().Where(models.User{Username: req.UserId}).Or(models.User{Email: req.UserId}).First(&user).Error; err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}
