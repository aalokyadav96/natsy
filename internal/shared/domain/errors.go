package domain

import "errors"

// Common domain sentinel errors across modules
var (
	ErrNotFound          = errors.New("requested resource was not found")
	ErrUnauthorized      = errors.New("unauthorized action")
	ErrInvalidInput      = errors.New("invalid input provided")
	ErrInternalError     = errors.New("internal server error")
	ErrConflict          = errors.New("resource conflict occurred")
	ErrInsufficientStock = errors.New("insufficient product stock")
)
