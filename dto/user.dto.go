package dto

import "github.com/fvrvz/auth-service-go/models"

type UserDTO struct {
	BaseModel
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	DOB       string `json:"DOB"`
	Username  string `json:"username"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=100"`
	FirstName string `json:"firstName" binding:"required,min=1,max=50"`
	LastName  string `json:"lastName" binding:"required,min=1,max=50"`
	DOB       string `json:"DOB" binding:"required,datetime=2006-01-02"`
}

func ToUserDTO(user *models.User) *UserDTO {
	return &UserDTO{
		FullName:  user.FirstName + " " + user.LastName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		DOB:       user.DOB.String(),
		BaseModel: BaseModel{
			CreatedAt:  user.CreatedAt,
			CreatedBy:  user.Username,
			ModifiedAt: &user.UpdatedAt,
			ModifiedBy: nil,
		},
	}
}
