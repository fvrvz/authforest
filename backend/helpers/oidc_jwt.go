package helpers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fvrvz/authforest/dto"
	"github.com/fvrvz/authforest/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// IDTokenClaims represents standard OIDC ID Token claims
type IDTokenClaims struct {
	jwt.RegisteredClaims
	Nonce             string `json:"nonce,omitempty"`
	AuthTime          int64  `json:"auth_time,omitempty"`
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
}

// AccessTokenClaims adds scope to the standard registered claims
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope,omitempty"`
}

// GenerateOIDCAccessToken creates an RS256-signed access token.
// userID should be the user's UUID (same as the sub claim in the ID token).
func GenerateOIDCAccessToken(cfg *dto.Config, userID string, clientID string, scopes []string, expiryMinutes int) (string, error) {
	now := time.Now()

	if expiryMinutes <= 0 {
		expiryMinutes = cfg.JWT.ExpiryMinutes
	}

	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			Issuer:    cfg.OIDC.Issuer,
			Audience:  jwt.ClaimStrings{clientID},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryMinutes) * time.Minute)),
		},
		Scope: strings.Join(scopes, " "),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = GetRSAKeyID()

	return token.SignedString(GetRSAPrivateKey())
}

// GenerateIDToken creates an RS256-signed OIDC ID Token.
// authTime is the time the user originally authenticated (per OIDC Core §2).
func GenerateIDToken(cfg *dto.Config, user *models.User, clientID string, nonce string, scopes []string, authTime time.Time, expiryMinutes int) (string, error) {
	now := time.Now()

	if expiryMinutes <= 0 {
		expiryMinutes = cfg.JWT.ExpiryMinutes
	}

	claims := IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   user.ID.String(),
			Issuer:    cfg.OIDC.Issuer,
			Audience:  jwt.ClaimStrings{clientID},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryMinutes) * time.Minute)),
		},
		Nonce:    nonce,
		AuthTime: authTime.Unix(),
	}

	// Include profile claims if "profile" scope requested
	if ContainsScope(scopes, "profile") {
		claims.Name = user.FirstName + " " + user.LastName
		claims.GivenName = user.FirstName
		claims.FamilyName = user.LastName
		claims.PreferredUsername = user.Username
	}

	// Include email claims if "email" scope requested
	if ContainsScope(scopes, "email") {
		claims.Email = user.Email
		claims.EmailVerified = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = GetRSAKeyID()

	return token.SignedString(GetRSAPrivateKey())
}

// GenerateOIDCRefreshToken creates an RS256-signed refresh token.
// userID should be the user's UUID for consistent sub claims.
func GenerateOIDCRefreshToken(cfg *dto.Config, userID string, clientID string, scopes []string, expiryHours int) (string, error) {
	now := time.Now()

	if expiryHours <= 0 {
		expiryHours = cfg.JWT.RefreshTokenExpiryHours
	}

	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			Issuer:    cfg.OIDC.Issuer,
			Audience:  jwt.ClaimStrings{clientID},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryHours) * time.Hour)),
		},
		Scope: strings.Join(scopes, " "),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = GetRSAKeyID()

	return token.SignedString(GetRSAPrivateKey())
}

// ExtractRSAClaims verifies and extracts standard claims from an RS256-signed token
func ExtractRSAClaims(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return GetRSAPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractOIDCAccessTokenClaims verifies and extracts access token claims including scope
func ExtractOIDCAccessTokenClaims(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return GetRSAPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ContainsScope checks if a scope list contains a specific scope value
func ContainsScope(scopes []string, target string) bool {
	for _, s := range scopes {
		if s == target {
			return true
		}
	}
	return false
}
