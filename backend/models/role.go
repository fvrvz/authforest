package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"default:''"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;primarykey"`
	RoleID uuid.UUID `gorm:"type:uuid;primarykey"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role   Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}
