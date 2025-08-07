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

	var user models.User

	if err := db.GetDB().Where("username = ? OR email = ?", req.UserId, req.UserId).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	config := config.GetConfig()
	token, refreshToken, err := helpers.GenerateJWT(config, user.Username)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	refreshTokenRecord := models.AuthRefreshTokens{
		Username:   user.Username,
		TokenHash:  helpers.Generate256Hash(refreshToken),
		ExpiryTime: time.Now().Add(time.Duration(config.JWT.RefreshTokenExpiryHours) * time.Hour),
	}

	if err := db.GetDB().Create(&refreshTokenRecord).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    constants.BEARER,
		ExpiresIn:    config.JWT.ExpiryMinutes * 60,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse[string]{
		Message: "User logged out successfully",
	})
}

func RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Invalid Input",
			Description: err.Error(),
		})
		return
	}

	claims, err := helpers.Verify(req.RefreshToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Refresh Token Expired",
			Description: err.Error(),
		})
		return
	}

	incomingTokenHash := helpers.Generate256Hash(req.RefreshToken)

	var record models.AuthRefreshTokens
	result := db.GetDB().Model(&models.AuthRefreshTokens{}).Where("token_hash = ? AND expiry_time > ?", incomingTokenHash, time.Now()).First(&record)

	if err := result.Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Database error",
			Description: err.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Refresh token expired",
		})
		return
	}

	deleteResult := db.GetDB().Model(&models.AuthRefreshTokens{}).Where("token_hash = ?", incomingTokenHash).Delete(&record)

	if err := deleteResult.Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Database error",
			Description: err.Error(),
		})
		return
	}

	if deleteResult.RowsAffected == 0 {
		log.Println("Token Expired")
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Error:       "Token Expired",
			Description: "Please login",
		})
		return
	}

	config := config.GetConfig()
	token, refreshToken, err := helpers.GenerateJWT(config, claims.Subject)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	refreshTokenRecord := models.AuthRefreshTokens{
		Username:   claims.Subject,
		TokenHash:  helpers.Generate256Hash(refreshToken),
		ExpiryTime: time.Now().Add(time.Duration(config.JWT.RefreshTokenExpiryHours) * time.Hour),
	}

	if err := db.GetDB().Create(&refreshTokenRecord).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to save refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    constants.BEARER,
		ExpiresIn:    config.JWT.ExpiryMinutes * 60,
	})
}
