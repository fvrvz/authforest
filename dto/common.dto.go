package dto

import "time"

type BaseModel struct {
	ID        uint64     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	CreatedBy string     `json:"createdBy"`
	UpdatedAt *time.Time `json:"updatedAt"`
	UpdatedBy *string    `json:"updatedBy"`
}

type SuccessResponse[TData any] struct {
	Data    TData  `json:"data,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}
