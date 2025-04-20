package middleware

import (
	"go-fiber-gorm/pkg/monitoring"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MonitorRequests returns a middleware that records request metrics
func MonitorRequests() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Route().Path // Get the route pattern, not the concrete URL
		method := c.Method()

		// Process request
		err := c.Next()

		// Record metrics after processing
		duration := time.Since(start).Seconds()
		statusCode := c.Response().StatusCode()

		// Record HTTP metrics
		monitoring.IncrementHTTPRequests(method, path, statusCode)
		monitoring.ObserveHTTPRequestLatency(method, path, duration)

		return err
	}
}
