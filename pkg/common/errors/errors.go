package errors

import (
	"fmt"
)

type AppError struct {
	StatusCode int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Field      string      `json:"field,omitempty"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func (e *AppError) WithField(field string) *AppError {
	e.Field = field
	return e
}

func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// Common error codes
const (
	CodeBadRequest          = "BAD_REQUEST"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeNotFound            = "NOT_FOUND"
	CodeConflict            = "CONFLICT"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeValidationError     = "VALIDATION_ERROR"
	CodeNotImplemented      = "NOT_IMPLEMENTED"
	CodeTooManyRequests     = "TOO_MANY_REQUESTS"
)
