package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/constants"
	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		TokenType:   constants.BEARER,
		ExpiresIn:   config.JWT.ExpiryHours * 3600,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse[string]{
		Message: "User logged out successfully",
	})
}

func Verify(ctx *gin.Context) (*jwt.RegisteredClaims, error) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != constants.BEARER {
		return nil, errors.New("authorization header format must be Bearer {token}")
	}

	tokenString := parts[1]

	jwtSecret := config.GetConfig().JWT.JWTSecret

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid Token")
	}

	return claims, nil
}
