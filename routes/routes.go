package routes

import (
	"go-fiber-gorm/config"
	"go-fiber-gorm/internal/handler"
	"go-fiber-gorm/internal/middleware"
	"go-fiber-gorm/pkg/docs"
	"go-fiber-gorm/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, cfg *config.Config, userHandler *handler.UserHandler, healthHandler *handler.HealthHandler) {
	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.MonitorRequests()) // Add monitoring middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// API Rate Limiter for all routes
	if cfg.Server.Env == "production" {
		app.Use(middleware.RateLimiter())
	}

	// Set up custom error handling middleware
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return errors.ErrorHandler(c, err)
		}
		return nil
	})

	// Health check endpoints
	app.Get("/health", healthHandler.Check)
	app.Get("/health/details", healthHandler.DetailedCheck)

	// Setup Swagger documentation
	docs.SetupSwagger(app, cfg)

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User routes
	users := v1.Group("/users")
	users.Post("/", middleware.ValidateRequest(handler.CreateUserRequest{}), userHandler.Create)
	users.Get("/", userHandler.GetAll)

	// Protected user routes
	usersProtected := users.Use(middleware.Protected())
	usersProtected.Get("/:id", userHandler.GetByID)
	usersProtected.Put("/:id", middleware.ValidateRequest(handler.UpdateUserRequest{}), userHandler.Update)
	usersProtected.Delete("/:id", userHandler.Delete)

	// Admin routes example
	admin := v1.Group("/admin").Use(middleware.Protected(), middleware.HasRole("admin"))
	admin.Get("/users", userHandler.GetAll) // Admin-only route example

	// Default 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Endpoint not found",
			},
		})
	})
}
