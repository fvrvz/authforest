package models

import "time"

type AuthRefreshTokens struct {
	JTI           string    `gorm:"primarykey;unique"`
	Username      string    `gorm:"not null"`
	ExpiresAt     time.Time `gorm:"not null"`
	IssuedAt      time.Time `gorm:"not null"`
	AccessTokenID string    `gorm:"not null"`
	Scope         string    `gorm:"default:''"` // space-separated scopes for OIDC tokens
}

type AccessTokenBlacklist struct {
	JTI       string    `gorm:"primarykey;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	IssuedAt  time.Time `gorm:"not null"`
}
