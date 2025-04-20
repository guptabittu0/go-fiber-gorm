package health

import (
	"github.com/gofiber/fiber/v2"
)

// Controller handles health check endpoints
type Controller struct {
	service *Service
}

// NewController creates a new health controller
func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

// RegisterRoutes registers the routes for the health module
func (c *Controller) RegisterRoutes(router fiber.Router) {
	health := router.Group("/health")

	health.Get("/", c.Check)
	health.Get("/details", c.DetailedCheck)
}

// Check handles basic health check
// @Summary      Health check endpoint
// @Description  Get application health status
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func (c *Controller) Check(ctx *fiber.Ctx) error {
	return ctx.JSON(c.service.CheckBasic())
}

// DetailedCheck handles detailed health check with all component statuses
// @Summary      Detailed health check
// @Description  Get detailed health status of all system components
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health/details [get]
func (c *Controller) DetailedCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(c.service.CheckDetailed())
}
