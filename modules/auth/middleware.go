package auth

import (
	"go-fiber-gorm/core/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware represents auth middleware functions
type Middleware struct {
	service *Service
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(service *Service) *Middleware {
	return &Middleware{
		service: service,
	}
}

// Protected ensures the request is authenticated
func (m *Middleware) Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return errors.NewUnauthorizedError("Authorization header is missing")
		}

		// Check the format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.NewUnauthorizedError("Authorization header format must be 'Bearer {token}'")
		}

		// Validate the token
		tokenString := parts[1]
		claims, err := m.service.ValidateToken(tokenString)
		if err != nil {
			return err
		}

		// Store user info in context
		ctx.Locals("userID", claims.UserID)
		ctx.Locals("userEmail", claims.Email)
		ctx.Locals("userRole", claims.Role)

		return ctx.Next()
	}
}

// RoleRequired ensures the authenticated user has the required role
func (m *Middleware) RoleRequired(role string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// First check if user is authenticated
		if err := m.Protected()(ctx); err != nil {
			return err
		}

		// Check if user has the required role
		userRole, ok := ctx.Locals("userRole").(string)
		if !ok || userRole != role {
			return errors.NewForbiddenError("You don't have permission to access this resource")
		}

		return ctx.Next()
	}
}

// GetAuthUser extracts the authenticated user from the context
func GetAuthUser(ctx *fiber.Ctx) (*Claims, error) {
	userID, ok1 := ctx.Locals("userID").(uint)
	email, ok2 := ctx.Locals("userEmail").(string)
	role, ok3 := ctx.Locals("userRole").(string)

	if !ok1 || !ok2 || !ok3 {
		return nil, errors.NewUnauthorizedError("User not authenticated")
	}

	return &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}
