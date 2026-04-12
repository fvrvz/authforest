package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OAuthClient struct {
	ID                       uuid.UUID      `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	ClientID                 string         `gorm:"unique;not null"`
	ClientSecret             string         `gorm:"not null"` // bcrypt hashed; empty for public clients
	ClientName               string         `gorm:"not null"`
	ClientType               string         `gorm:"not null;default:'public'"` // "public" or "confidential"
	RedirectURIs             pq.StringArray `gorm:"type:text[];not null"`
	Scopes                   string         `gorm:"not null;default:'openid profile email'"` // space-separated
	GrantTypes               string         `gorm:"not null;default:'authorization_code'"`   // space-separated
	AccessTokenExpiryMinutes int            `gorm:"not null;default:15"`
	RefreshTokenExpiryHours  int            `gorm:"not null;default:2"`
	IDTokenExpiryMinutes     int            `gorm:"not null;default:15"`
	CreatedAt                time.Time
	UpdatedAt                *time.Time
}

type AuthorizationCode struct {
	Code                string `gorm:"primarykey"`
	ClientID            string `gorm:"not null;index"`
	Username            string `gorm:"not null"`
	RedirectURI         string `gorm:"not null"`
	Scope               string `gorm:"not null"`
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time `gorm:"not null"`
	CreatedAt           time.Time
	Used                bool `gorm:"not null;default:false"`
}
