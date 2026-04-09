package dto

type RoleDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
