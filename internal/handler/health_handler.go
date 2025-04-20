package handler

import (
	"context"
	"go-fiber-gorm/internal/repository"
	"go-fiber-gorm/pkg/cache"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles basic health check
// @Summary      Health check endpoint
// @Description  Get application health status
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"service": "fiber-gorm-api",
	})
}

// DetailedCheck handles detailed health check with all component statuses
// @Summary      Detailed health check
// @Description  Get detailed health status of all system components
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health/details [get]
func (h *HealthHandler) DetailedCheck(c *fiber.Ctx) error {
	// Check database connection
	dbStatus := "healthy"
	sqlDB, err := repository.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "unhealthy"
	}

	// Check Redis connection if available
	redisStatus := "not configured"
	if cache.Client != nil {
		redisStatus = "healthy"
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err := cache.Client.Ping(ctx).Result()
		if err != nil {
			redisStatus = "unhealthy"
		}
	}

	// Get system info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return c.JSON(fiber.Map{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"service": "fiber-gorm-api",
		"components": fiber.Map{
			"database": dbStatus,
			"redis":    redisStatus,
		},
		"system": fiber.Map{
			"memory": fiber.Map{
				"alloc":      m.Alloc / 1024 / 1024,
				"totalAlloc": m.TotalAlloc / 1024 / 1024,
				"sys":        m.Sys / 1024 / 1024,
				"numGC":      m.NumGC,
			},
			"goroutines": runtime.NumGoroutine(),
		},
	})
}
