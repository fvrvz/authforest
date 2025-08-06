package services

import (
	"net/http"

	"github.com/fvrvz/auth-service-go/config"
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
	token, err := helpers.GenerateJWT(config, user.Username, user.Email)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   config.JWT.ExpiryHours * 3600,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse[string]{
		Message: "User logged out successfully",
	})
}
