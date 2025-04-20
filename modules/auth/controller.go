package auth

import (
	"go-fiber-gorm/core/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Controller handles auth-related HTTP requests
type Controller struct {
	service *Service
}

// NewController creates a new auth controller
func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

// RegisterRoutes registers all auth-related routes
func (c *Controller) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/register", c.Register)
	auth.Post("/login", c.Login)
	auth.Post("/refresh-token", c.RefreshToken)

	// Protected routes
	auth.Post("/logout", c.AuthMiddleware(), c.Logout)
	auth.Post("/logout-all", c.AuthMiddleware(), c.LogoutAll)
	auth.Post("/change-password", c.AuthMiddleware(), c.ChangePassword)
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (c *Controller) Register(ctx *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return errors.NewBadRequestError("Invalid request body")
	}

	result, err := c.service.Register(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// Login handles user login
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (c *Controller) Login(ctx *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return errors.NewBadRequestError("Invalid request body")
	}

	result, err := c.service.Login(req)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh-token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/refresh-token [post]
func (c *Controller) RefreshToken(ctx *fiber.Ctx) error {
	req := new(RefreshTokenRequest)
	if err := ctx.BodyParser(req); err != nil {
		return errors.NewBadRequestError("Invalid request body")
	}

	result, err := c.service.RefreshToken(req)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate the current session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/logout [post]
func (c *Controller) Logout(ctx *fiber.Ctx) error {
	// Get refresh token from request
	req := new(RefreshTokenRequest)
	if err := ctx.BodyParser(req); err != nil || req.RefreshToken == "" {
		return errors.NewBadRequestError("Refresh token is required")
	}

	if err := c.service.Logout(req.RefreshToken); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// LogoutAll handles logging out all sessions
// @Summary Logout all sessions
// @Description Invalidate all sessions for the current user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/logout-all [post]
func (c *Controller) LogoutAll(ctx *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return errors.NewUnauthorizedError("User not authenticated")
	}

	if err := c.service.LogoutAll(userID); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Logged out from all devices successfully",
	})
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change the user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body ChangePasswordRequest true "Password change request"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/change-password [post]
func (c *Controller) ChangePassword(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return errors.NewUnauthorizedError("User not authenticated")
	}

	// Parse request
	req := new(ChangePasswordRequest)
	if err := ctx.BodyParser(req); err != nil {
		return errors.NewBadRequestError("Invalid request body")
	}

	if err := c.service.ChangePassword(userID, req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}

// AuthMiddleware returns a middleware that checks authentication
func (c *Controller) AuthMiddleware() fiber.Handler {
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
		claims, err := c.service.ValidateToken(tokenString)
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
