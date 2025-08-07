package services

import (
	"fmt"
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

	accessToken, refreshToken, err := generateAndReplaceRefreshToken(user.Username)

	if err != nil {
		log.Printf("Login token error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Token generation failed"})
		return
	}

	respondWithTokens(c, accessToken, refreshToken)
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
			Error:       "Invalid or expired refresh token",
			Description: err.Error(),
		})
		return
	}

	record, err := findValidRefreshToken(req.RefreshToken)

	if err != nil {
		log.Printf("Refresh token lookup failed: %v", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Refresh token not valid"})
		return
	}

	if err := deleteRefreshToken(record); err != nil {
		log.Printf("Failed to delete refresh token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
		return
	}

	accessToken, newRefreshToken, err := generateAndStoreNewTokens(claims.Subject)

	if err != nil {
		log.Printf("Failed to generate new tokens: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Token generation failed",
			Description: err.Error(),
		})
		return
	}

	respondWithTokens(ctx, accessToken, newRefreshToken)
}

func findValidRefreshToken(refreshToken string) (*models.AuthRefreshTokens, error) {
	hashedToken := helpers.Generate256Hash(refreshToken)

	var record *models.AuthRefreshTokens

	if err := db.GetDB().Where("token_hash = ? AND expiry_time > ?", hashedToken, time.Now()).First(&record).Error; err != nil {
		return nil, err
	}

	return record, nil
}

func deleteRefreshToken(record *models.AuthRefreshTokens) error {
	return db.GetDB().Delete(&record).Error
}

func generateAndStoreNewTokens(username string) (string, string, error) {
	config := config.GetConfig()

	accessToken, refreshToken, err := helpers.GenerateJWT(config, username)

	if err != nil {
		return "", "", err
	}

	newRecord := models.AuthRefreshTokens{
		Username:   username,
		TokenHash:  helpers.Generate256Hash(refreshToken),
		ExpiryTime: time.Now().Add(time.Duration(config.JWT.RefreshTokenExpiryHours) * time.Hour),
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

	if err := db.GetDB().Where("username = ? OR email = ?", req.UserId, req.UserId).First(&user).Error; err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func generateAndReplaceRefreshToken(username string) (string, string, error) {
	config := config.GetConfig()
	accessToken, refreshToken, err := helpers.GenerateJWT(config, username)

	if err != nil {
		return "", "", err
	}

	var existing models.AuthRefreshTokens

	if err := db.GetDB().Where("username = ?", username).First(&existing).Error; err == nil {
		if err := deleteRefreshToken(&existing); err != nil {
			return "", "", fmt.Errorf("failed to delete existing refresh token: %w", err)
		}
	}

	newRecord := models.AuthRefreshTokens{
		Username:   username,
		TokenHash:  helpers.Generate256Hash(refreshToken),
		ExpiryTime: time.Now().Add(time.Duration(config.JWT.RefreshTokenExpiryHours) * time.Hour),
	}

	if err := db.GetDB().Create(&newRecord).Error; err != nil {
		return "", "", fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
