package user

import (
	"go-fiber-gorm/core/errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Controller handles HTTP requests related to users
type Controller struct {
	service *Service
}

// NewController creates a new user controller
func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

// RegisterRoutes registers the routes for the user module
func (c *Controller) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users")

	// Public routes
	users.Post("/", c.Create)
	users.Get("/", c.GetAll)

	// Protected routes - in a real app, apply auth middleware here
	users.Get("/:id", c.GetByID)
	users.Put("/:id", c.Update)
	users.Delete("/:id", c.Delete)
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
func (c *Controller) Create(ctx *fiber.Ctx) error {
	req := new(CreateUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		return errors.NewBadRequestError("Invalid request body")
	}

	user, err := c.service.Create(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
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
func (c *Controller) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return errors.NewBadRequestError("Invalid user ID")
	}

	user, err := c.service.GetByID(uint(id))
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
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
func (c *Controller) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return errors.NewBadRequestError("Invalid user ID")
	}

	req := new(UpdateUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		return errors.NewBadRequestError("Invalid request body")
	}

	user, err := c.service.Update(uint(id), req)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
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
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return errors.NewBadRequestError("Invalid user ID")
	}

	if err := c.service.Delete(uint(id)); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
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
func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	// Parse query parameters for pagination
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	users, count, err := c.service.GetAll(page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
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
