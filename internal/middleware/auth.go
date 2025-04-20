package middleware

import (
	"go-fiber-gorm/pkg/auth"
	"go-fiber-gorm/pkg/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTSecret stores the JWT secret key
var JWTSecret string

// SetJWTSecret sets the JWT secret key for auth middleware
func SetJWTSecret(secret string) {
	JWTSecret = secret
}

// Protected protects routes that require authentication
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")

		// Check if auth header exists and has right format
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return errors.NewUnauthorizedError("Missing or invalid authorization token")
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		userID, err := auth.GetUserIDFromToken(tokenString, JWTSecret)
		if err != nil {
			return err
		}

		// Set user ID in context for subsequent middleware/handlers
		c.Locals("userID", userID)

		// Get user role (optional)
		role, err := auth.GetUserRoleFromToken(tokenString, JWTSecret)
		if err == nil && role != "" {
			c.Locals("userRole", role)
		}

		return c.Next()
	}
}

// HasRole checks if the user has the required role
func HasRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok || userRole != role {
			return errors.NewForbiddenError("Insufficient permissions")
		}

		return c.Next()
	}
}
