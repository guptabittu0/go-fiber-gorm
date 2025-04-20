package errors

import (
	"errors"
	"fmt"
	"go-fiber-gorm/core/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// AppError represents an application error
type AppError struct {
	StatusCode int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

// Error returns the error message
func (e *AppError) Error() string {
	return e.Message
}

// ErrorHandler handles all errors in the application
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Default error
	code := "INTERNAL_SERVER_ERROR"
	message := "Something went wrong"
	statusCode := fiber.StatusInternalServerError
	details := make(map[string]interface{})

	// Check if it's our custom AppError
	var appError *AppError
	if errors.As(err, &appError) {
		code = appError.Code
		message = appError.Message
		statusCode = appError.StatusCode
		if detailsMap, ok := appError.Details.(map[string]interface{}); ok {
			details = detailsMap
		} else if appError.Details != nil {
			// If it's not a map but contains data, add it as a value
			details["error_details"] = appError.Details
		}

		// Only log server errors
		if statusCode >= 500 {
			logger.Error(fmt.Sprintf("[%s] %s", code, message))
		}
	} else {
		// Handle validation errors
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			code = "VALIDATION_ERROR"
			message = "Validation failed"
			statusCode = fiber.StatusBadRequest

			// Format validation errors
			for _, err := range validationErrors {
				field := err.Field()
				tag := err.Tag()
				details[field] = fmt.Sprintf("Failed on '%s' validation", tag)
			}
		} else {
			// This is an unexpected error
			logger.Error(fmt.Sprintf("Unexpected error: %v", err))
		}
	}

	// Return error response
	return ctx.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}

// New creates a new AppError
func New(statusCode int, code string, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// Common error constructors
func NewBadRequestError(message string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized"
	}
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Forbidden"
	}
	return New(http.StatusForbidden, "FORBIDDEN", message)
}

func NewNotFoundError(resource string) *AppError {
	message := resource + " not found"
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

func NewValidationError(err error) *AppError {
	appError := New(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")

	// Parse validation errors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make(map[string]string)
		for _, err := range validationErrors {
			field := err.Field()
			tag := err.Tag()
			details[field] = fmt.Sprintf("Failed on '%s' validation", tag)
		}
		appError.Details = details
	} else {
		appError.Message = err.Error()
	}

	return appError
}

func NewInternalServerError(message string) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func NewTooManyRequestsError(message string) *AppError {
	if message == "" {
		message = "Too many requests"
	}
	return New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message)
}

func NewServiceUnavailableError(message string) *AppError {
	if message == "" {
		message = "Service unavailable"
	}
	return New(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}
