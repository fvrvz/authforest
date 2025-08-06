package helpers

import (
	"auth-service-go/initializers"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func GenerateJWT(cfg *initializers.Config, username string, email string) (string, error) {
	expiration := time.Now().Add(time.Duration(cfg.JWT.ExpiryHours) * time.Hour)

	claims := jwt.MapClaims{
		"sub":   username,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.JWT.JWTSecret))
}
