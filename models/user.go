package models

import (
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        uint64
	UpdatedBy *string
	CreatedBy string

	Username  string    `gorm:"primarykey;unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"unique;not null"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	DOB       time.Time `gorm:"column:dob"`
}
