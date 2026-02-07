package dto

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt" gorm:"-"`
	CreatedBy string     `json:"createdBy" gorm:"-"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"-"`
	UpdatedBy *string    `json:"updatedBy" gorm:"-"`
}

type SuccessResponse[TData any] struct {
	Data    TData  `json:"data,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}
