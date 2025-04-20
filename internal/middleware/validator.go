package middleware

import (
	"go-fiber-gorm/pkg/errors"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateRequest validates a struct based on validation tags
func ValidateRequest(model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new instance of the model
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}

		modelValue := reflect.New(modelType).Interface()

		// Parse request body into model
		if err := c.BodyParser(modelValue); err != nil {
			return errors.NewBadRequestError("Invalid request body")
		}

		// Validate the model
		if err := validate.Struct(modelValue); err != nil {
			return errors.NewValidationError(err)
		}

		// Store validated model in context for handlers to use
		c.Locals("validated", modelValue)

		// Continue to next middleware/handler
		return c.Next()
	}
}

// GetValidatedBody retrieves the validated request body from context
func GetValidatedBody(c *fiber.Ctx) interface{} {
	return c.Locals("validated")
}
