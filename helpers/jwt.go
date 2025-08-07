package helpers

import (
	"time"

	"github.com/fvrvz/auth-service-go/dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func GenerateJWT(cfg *dto.Config, username string) (accessToken string, refreshToken string, err error) {
	expiration := time.Now().Add(time.Duration(cfg.JWT.ExpiryMinutes) * time.Minute)

	accessClaims := jwt.RegisteredClaims{
		Subject:   username,
		Issuer:    "auth-service-go",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiration),
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessToken, err = access.SignedString([]byte(cfg.JWT.JWTSecret))

	if err != nil {
		return
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   username,
		Issuer:    "auth-service-go",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.RefreshTokenExpiryHours) * time.Hour)),
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshToken, err = refresh.SignedString([]byte(cfg.JWT.JWTSecret))

	if err != nil {
		return
	}

	return
}
