package dto

import "time"

type UserDTO struct {
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	DOB       string    `json:"DOB"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=100"`
	FirstName string `json:"firstName" binding:"required,min=1,max=50"`
	LastName  string `json:"lastName" binding:"required,min=1,max=50"`
	DOB       string `json:"DOB" binding:"required,datetime=2006-01-02"`
}
