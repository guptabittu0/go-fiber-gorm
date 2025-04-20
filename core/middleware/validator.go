package middleware

import (
	"go-fiber-gorm/core/errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ValidatorMiddleware creates a middleware for request validation
func ValidatorMiddleware() fiber.Handler {
	validate := validator.New()

	// Register custom validation methods if needed
	// validate.RegisterValidation("custom", customValidationFunc)

	// Register function to get tag names from json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return func(c *fiber.Ctx) error {
		// Check if the route has validation requirements
		validationStruct := c.Locals("validation_struct")
		if validationStruct == nil {
			return c.Next()
		}

		// Create a new instance of the struct type
		structType := reflect.TypeOf(validationStruct)
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
		}

		instance := reflect.New(structType).Interface()

		// Parse request body into the struct
		if err := c.BodyParser(instance); err != nil {
			return errors.NewBadRequestError("Invalid request body")
		}

		// Validate the struct
		if err := validate.Struct(instance); err != nil {
			return errors.NewValidationError(err)
		}

		// Store validated struct in context
		c.Locals("validated", instance)

		return c.Next()
	}
}

// Validate is a helper function to validate a struct in handlers
func Validate(validate *validator.Validate, obj interface{}) error {
	if err := validate.Struct(obj); err != nil {
		return errors.NewValidationError(err)
	}
	return nil
}
