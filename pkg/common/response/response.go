package response

import (
	"net/http"

	"golang-template/pkg/common/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorDetail struct {
	Code    string      `json:"code" example:"INTERNAL_SERVER_ERROR"`
	Message string      `json:"message" example:"An unexpected error occurred"`
	Field   string      `json:"field,omitempty" example:"email"`
	Details interface{} `json:"details,omitempty"`
}

type MetaData struct {
	Total       int64 `json:"total"`
	Count       int   `json:"count"`
	PerPage     int   `json:"perPage"`
	CurrentPage int   `json:"currentPage"`
	TotalPages  int   `json:"totalPages"`
	NextPage    *int  `json:"nextPage,omitempty"`
	PrevPage    *int  `json:"prevPage,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	var statusCode int
	var errorResponse interface{}

	if appErr, ok := err.(*errors.AppError); ok {
		statusCode = appErr.StatusCode
		errorResponse = ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
			Field:   appErr.Field,
			Details: appErr.Details,
		}
	} else {
		statusCode = http.StatusInternalServerError
		errorResponse = ErrorDetail{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred",
		}
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorResponse,
	})
}

// ErrorWithCode sends an error response with a custom status code
func ErrorWithCode(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: ErrorDetail{
			Code:    errors.CodeValidationError,
			Message: message,
			Field:   field,
		},
	})
}

// WithMeta sends a successful response with metadata
func WithMeta(c *gin.Context, statusCode int, data interface{}, meta MetaData) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// 200 OK response
func OK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, data)
}

// 201 Created response
func Created(c *gin.Context, data interface{}) {
	Success(c, http.StatusCreated, data)
}

// 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	ErrorWithCode(c, http.StatusBadRequest, errors.CodeBadRequest, message)
}

// 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	ErrorWithCode(c, http.StatusUnauthorized, errors.CodeUnauthorized, message)
}

// 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden access"
	}
	ErrorWithCode(c, http.StatusForbidden, errors.CodeForbidden, message)
}

// 404 Not Found response
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorWithCode(c, http.StatusNotFound, errors.CodeNotFound, message)
}

// 429 Too Many Requests response
func RateLimitExceeded(c *gin.Context) {
	ErrorWithCode(c, http.StatusTooManyRequests, errors.CodeTooManyRequests, "Rate limit exceeded")
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorWithCode(c, http.StatusInternalServerError, errors.CodeInternalServerError, message)
}

// creates pagination metadata from total count and pagination parameters
func CreatePaginationMetadata(total int64, page, limit int) MetaData {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	var nextPage *int
	var prevPage *int

	if page < totalPages {
		next := page + 1
		nextPage = &next
	}

	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	return MetaData{
		Total:       total,
		Count:       min(limit, int(total)-(page-1)*limit),
		PerPage:     limit,
		CurrentPage: page,
		TotalPages:  totalPages,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
