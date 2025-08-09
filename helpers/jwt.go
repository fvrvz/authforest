package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/constants"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func GenerateJWT(cfg *dto.Config, username string) (accessToken string, refreshToken string, err error) {
	now := time.Now()

	accessClaims := jwt.RegisteredClaims{
		ID:        uuid.New().String(),
		Subject:   username,
		Issuer:    "auth-service-go",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.ExpiryMinutes) * time.Minute)),
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessToken, err = access.SignedString([]byte(cfg.JWT.JWTSecret))

	if err != nil {
		return
	}

	refreshClaims := jwt.RegisteredClaims{
		ID:        uuid.New().String(),
		Subject:   username,
		Issuer:    "auth-service-go",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.RefreshTokenExpiryHours) * time.Hour)),
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshToken, err = refresh.SignedString([]byte(cfg.JWT.JWTSecret))

	if err != nil {
		return
	}

	return
}

func ExtractTokenFromHeaders(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != constants.BEARER {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

func ExtractClaims(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().JWT.JWTSecret), nil
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

func Generate256Hash(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}
