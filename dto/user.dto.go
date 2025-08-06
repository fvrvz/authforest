package dto

import "time"

type UserDTO struct {
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	DOB       time.Time `json:"DOB"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
