package main

import (
	"fmt"
	"go-fiber-gorm/config"
	"go-fiber-gorm/core/cache"
	"go-fiber-gorm/core/database"
	appErrors "go-fiber-gorm/core/errors"
	"go-fiber-gorm/core/logger"
	"go-fiber-gorm/core/middleware"
	"go-fiber-gorm/migrations"
	"go-fiber-gorm/modules/auth"
	"go-fiber-gorm/modules/user"
	"go-fiber-gorm/routes"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize environment
	if err := config.LoadEnvForCurrentEnvironment(); err != nil {
		fmt.Printf("Failed to load environment: %v\n", err)
		os.Exit(1)
	}

	// Create .env file from example if it doesn't exist
	if err := config.CreateEnvFile(); err != nil {
		fmt.Printf("Failed to create .env file: %v\n", err)
		// Not fatal, continue with defaults
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	logger.Setup(cfg.Server.Env)
	logger.Info("Starting application...")
	logger.Info("Environment:", cfg.Server.Env)
	logger.Info("Go Version:", runtime.Version())

	// Connect to database
	dbConn, err := database.NewConnection(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	// Get the GORM DB instance
	db := dbConn.GetDB()

	// Connect to Redis (optional - will continue if Redis isn't available)
	redisClient, err := cache.ConnectRedis(&cfg.Redis)
	if err != nil {
		logger.Warn("Failed to connect to Redis, continuing without cache:", err)
	}

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations:", err)
	}

	// ?? Initialize worker pool for background tasks
	// workerPool := worker.NewPool(runtime.NumCPU())
	// workerPool.Start()
	// defer workerPool.Stop()

	// // Submit example background task
	// workerPool.Submit(func() error {
	// 	logger.Info("Running startup tasks in background worker...")
	// 	// Your background task logic here
	// 	return nil
	// })

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		Prefork:           cfg.Server.Env == "production",
		EnablePrintRoutes: cfg.Server.Env != "production",
		ErrorHandler:      appErrors.ErrorHandler,
	})

	// Apply global middleware
	app.Use(fiberLogger.New())
	// app.Use(middleware.Logger())
	app.Use(middleware.RateLimiter())

	// Register models for auto-migration
	if err := dbConn.AutoMigrate(
		&user.User{},
		&auth.Session{},
	); err != nil {
		logger.Fatal("Failed to auto migrate models:", err)
	}

	// Setup routes using the new modular structure
	routes.SetupRoutes(app, cfg, db, redisClient)

	// Start server in a goroutine
	go func() {
		serverPort := cfg.Server.Port
		logger.Info("Server starting on port ", serverPort)
		if err := app.Listen(":" + serverPort); err != nil {
			logger.Fatal("Failed to start server:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Close Redis connection if it was established
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			logger.Error("Error closing Redis connection:", err)
		}
		logger.Info("Redis connection closed")
	}

	// Shutdown fiber app
	if err := app.Shutdown(); err != nil {
		logger.Fatal("Server shutdown failed:", err)
	}
	logger.Info("Server gracefully stopped")
}
