package dto

import "time"

type BaseModel struct {
	ID         string     `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	CreatedBy  string     `json:"createdBy"`
	ModifiedAt *time.Time `json:"modifiedAt"`
	ModifiedBy *string    `json:"modifiedBy"`
}

type SuccessResponse[TData any] struct {
	Data    TData  `json:"data,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}
