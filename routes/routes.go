package routes

import (
	"go-fiber-gorm/modules/auth"
	"go-fiber-gorm/modules/health"
	"go-fiber-gorm/modules/user"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupRoutes configures the application routes and middleware
func SetupRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	// Global middleware
	app.Use(cors.New())
	app.Use(recover.New())

	// API routes with version prefix
	api := app.Group("/api/v1")

	// Health module setup
	healthService := health.NewService(db, redisClient) // Replace nil with redis client if available
	healthController := health.NewController(healthService)
	healthController.RegisterRoutes(api)

	// User module setup
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userController := user.NewController(userService)

	// Auth module setup
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(
		authRepo,
		userRepo,
		auth.ServiceConfig{
			JWTSecret:     "your-secret-key",  // Should be loaded from config
			AccessExpiry:  time.Hour * 1,      // 1 hour
			RefreshExpiry: time.Hour * 24 * 7, // 7 days
		},
	)
	authMiddleware := auth.NewMiddleware(authService)
	authController := auth.NewController(authService)

	// Register auth routes
	authController.RegisterRoutes(api)

	// Register user routes (using auth middleware for protected routes)
	users := api.Group("/users")
	users.Post("/", authMiddleware.RoleRequired("admin"), userController.Create)
	users.Get("/", userController.GetAll)
	users.Get("/:id", authMiddleware.Protected(), userController.GetByID)
	users.Put("/:id", authMiddleware.Protected(), userController.Update)
	users.Delete("/:id", authMiddleware.RoleRequired("admin"), userController.Delete)

	// 404 Handler
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
