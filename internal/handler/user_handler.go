package handler

import (
	"strconv"

	"go-fiber-gorm/internal/model"
	"go-fiber-gorm/internal/service"
	"go-fiber-gorm/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// Request/response types
type (
	// CreateUserRequest represents the request to create a user
	CreateUserRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	// UpdateUserRequest represents the request to update a user
	UpdateUserRequest struct {
		Name  string `json:"name,omitempty" validate:"omitempty,min=2"`
		Email string `json:"email,omitempty" validate:"omitempty,email"`
	}

	// UserResponse represents the user response
	UserResponse struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	// UsersResponse represents a paginated list of users
	UsersResponse struct {
		Users []UserResponse `json:"users"`
		Meta  PaginationMeta `json:"meta"`
	}

	// PaginationMeta represents pagination metadata
	PaginationMeta struct {
		Total int64 `json:"total"`
		Page  int   `json:"page"`
		Limit int   `json:"limit"`
		Pages int64 `json:"pages"`
	}
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Create handles user creation
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User information"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	// Get validated request from context (set by ValidateRequest middleware)
	req, ok := c.Locals("validated").(*CreateUserRequest)
	if !ok {
		// If not available, try parsing the body directly
		req = new(CreateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return errors.NewBadRequestError("Invalid request body")
		}
	}

	// Map to model request
	modelReq := &model.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service
	user, err := h.userService.Create(modelReq)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// GetByID handles retrieving a user by ID
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return errors.NewBadRequestError("Invalid user ID")
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// Update handles updating a user
// @Summary Update a user
// @Description Update a user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return errors.NewBadRequestError("Invalid user ID")
	}

	// Get validated request from context
	req, ok := c.Locals("validated").(*UpdateUserRequest)
	if !ok {
		req = new(UpdateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return errors.NewBadRequestError("Invalid request body")
		}
	}

	// Map to model request
	modelReq := &model.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := h.userService.Update(uint(id), modelReq)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// Delete handles deleting a user
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return errors.NewBadRequestError("Invalid user ID")
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

// GetAll handles retrieving all users with pagination
// @Summary Get all users
// @Description Get all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users [get]
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	// Parse query parameters for pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	users, count, err := h.userService.GetAll(page, limit)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"users": users,
			"meta": fiber.Map{
				"total": count,
				"page":  page,
				"limit": limit,
				"pages": (count + int64(limit) - 1) / int64(limit),
			},
		},
	})
}
