package dto

import "github.com/fvrvz/auth-service-go/models"

type UserDTO struct {
	BaseModel
	FullName  string    `json:"fullName" gorm:"->"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	DOB       string    `json:"DOB"`
	Username  string    `json:"username"`
	Roles     []RoleDTO `json:"roles" gorm:"-"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=100"`
	FirstName string `json:"firstName" binding:"required,min=1,max=50"`
	LastName  string `json:"lastName" binding:"required,min=1,max=50"`
	DOB       string `json:"DOB" binding:"required,datetime=2006-01-02"`
}

type AdminCreateUserRequest struct {
	Username  string   `json:"username" binding:"required,min=3,max=20"`
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required,min=6,max=100"`
	FirstName string   `json:"firstName" binding:"required,min=1,max=50"`
	LastName  string   `json:"lastName" binding:"required,min=1,max=50"`
	DOB       string   `json:"DOB" binding:"required,datetime=2006-01-02"`
	RoleIDs   []string `json:"roleIds"`
}

type UpdateUserRequest struct {
	Email     *string `json:"email" binding:"omitempty,email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	DOB       *string `json:"DOB" binding:"omitempty,datetime=2006-01-02"`
}

type AdminUpdateUserRequest struct {
	Email     *string  `json:"email" binding:"omitempty,email"`
	FirstName *string  `json:"firstName"`
	LastName  *string  `json:"lastName"`
	DOB       *string  `json:"DOB" binding:"omitempty,datetime=2006-01-02"`
	RoleIDs   []string `json:"roleIds"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

func ToUserDTO(user *models.User) *UserDTO {
	roles := make([]RoleDTO, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = RoleDTO{
			ID:          r.ID.String(),
			Name:        r.Name,
			Description: r.Description,
		}
	}
	return &UserDTO{
		FullName:  user.FirstName + " " + user.LastName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		DOB:       user.DOB.String(),
		Roles:     roles,
		BaseModel: BaseModel{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			CreatedBy: user.Username,
			UpdatedAt: user.UpdatedAt,
			UpdatedBy: user.UpdatedBy,
		},
	}
}
