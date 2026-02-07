package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	UpdatedBy *string
	CreatedBy string

	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"unique;not null"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	DOB       time.Time `gorm:"column:dob"`
}
