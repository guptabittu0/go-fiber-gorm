package middleware

import (
	"go-fiber-gorm/core/errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter creates a middleware that limits repeated requests to public APIs
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // Max number of requests within duration
		Expiration: 1 * time.Minute, // Duration for max requests
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Use IP as the rate limit key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return errors.NewTooManyRequestsError("Rate limit exceeded. Please try again later.")
		},
		Storage: nil, // Default in-memory storage
	})
}
