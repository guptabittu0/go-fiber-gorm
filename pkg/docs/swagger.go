package docs

import (
	"go-fiber-gorm/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Fiber Gorm API
// @version 1.0
// @description A production-ready API boilerplate using Go Fiber and GORM.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// SetupSwagger sets up Swagger documentation
func SetupSwagger(app *fiber.App, cfg *config.Config) {
	// Swagger endpoint
	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title:        "Fiber Gorm API",
		DeepLinking:  true,
		DocExpansion: "list",
	}))
}
