package middleware

import (
	"fmt"
	"go-fiber-gorm/core/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger returns a middleware which logs HTTP requests/responses
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		method := c.Method()

		// Handle the request
		err := c.Next()

		// Log after the request is done
		latency := time.Since(start)
		statusCode := c.Response().StatusCode()
		ip := c.IP()

		// Color codes for terminal output
		methodColor := "\033[36m" // Cyan for method
		pathColor := "\033[34m"   // Blue for path
		resetColor := "\033[0m"   // Reset color
		statusColor := "\033[32m" // Green for success (default)

		// Determine log level based on status code
		var logFn func(args ...interface{})
		if statusCode >= 500 {
			statusColor = "\033[31m" // Red for server errors
			logFn = logger.Error
		} else if statusCode >= 400 {
			statusColor = "\033[33m" // Yellow for client errors
			logFn = logger.Warn
		} else {
			logFn = logger.Info
		}

		// Use colors for different parts of the log message
		logFn(fmt.Sprintf("%s%s%s %s%s%s%s %d%s %s%s%s, IP %s",
			methodColor, method, resetColor,
			pathColor, path, resetColor,
			statusColor, statusCode, resetColor, fiber.DefaultColors.Yellow, latency, resetColor, ip))

		return err
	}
}
