package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

// AppError represents an application error
type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// Error returns the error message
func (e AppError) Error() string {
	return e.Message
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(entity string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", entity),
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

// NewTooManyRequestsError creates a new rate limit exceeded error
func NewTooManyRequestsError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusTooManyRequests,
		Code:       "TOO_MANY_REQUESTS",
		Message:    message,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
	}
}

// NewValidationError creates a new validation error from validator errors
func NewValidationError(err error) *AppError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, formatValidationError(e))
		}
		return &AppError{
			StatusCode: http.StatusBadRequest,
			Code:       "VALIDATION_ERROR",
			Message:    strings.Join(errorMessages, ", "),
		}
	}

	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    err.Error(),
	}
}

// formatValidationError formats a validation error into a human-readable message
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// ErrorHandler is a middleware to handle errors
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error
	statusCode := fiber.StatusInternalServerError
	code := "INTERNAL_SERVER_ERROR"
	message := "Something went wrong"

	// Check if it's an AppError
	if appErr, ok := err.(*AppError); ok {
		statusCode = appErr.StatusCode
		code = appErr.Code
		message = appErr.Message
	} else if fiberErr, ok := err.(*fiber.Error); ok {
		// Check if it's a Fiber error
		statusCode = fiberErr.Code
		message = fiberErr.Message
	} else {
		// For any other error, set generic error response
		message = err.Error()
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
