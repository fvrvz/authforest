package dto

type SuccessResponse[TData any] struct {
	Data    TData  `json:"data,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}
