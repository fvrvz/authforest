package models

import "time"

type AuthRefreshTokens struct {
	Username   string    `gorm:"primarykey;unique"`
	TokenHash  string    `gorm:"unique;not null;column:token_hash"`
	ExpiryTime time.Time `gorm:"not null"`
	CreatedAt  time.Time
}
