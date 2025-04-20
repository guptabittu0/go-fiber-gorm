package health

import (
	"context"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

// Service handles health check business logic
type Service struct {
	db        *gorm.DB
	redisConn *redis.Client
}

// NewService creates a new health service
func NewService(db *gorm.DB, redisConn *redis.Client) *Service {
	return &Service{
		db:        db,
		redisConn: redisConn,
	}
}

// CheckBasic performs a basic health check
func (s *Service) CheckBasic() map[string]interface{} {
	return map[string]interface{}{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"service": "fiber-gorm-api",
	}
}

// CheckDetailed performs a detailed health check of all components
func (s *Service) CheckDetailed() map[string]interface{} {
	// Check database connection
	dbStatus := "healthy"
	sqlDB, err := s.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "unhealthy"
	}

	// Check Redis connection if available
	redisStatus := "not configured"
	if s.redisConn != nil {
		redisStatus = "healthy"
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := s.redisConn.Ping(ctx).Err(); err != nil {
			redisStatus = "unhealthy"
		}
	}

	// Get system info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"service": "fiber-gorm-api",
		"components": map[string]interface{}{
			"database": dbStatus,
			"redis":    redisStatus,
		},
		"system": map[string]interface{}{
			"memory": map[string]interface{}{
				"alloc":      m.Alloc / 1024 / 1024,
				"totalAlloc": m.TotalAlloc / 1024 / 1024,
				"sys":        m.Sys / 1024 / 1024,
				"numGC":      m.NumGC,
			},
			"goroutines": runtime.NumGoroutine(),
		},
	}
}
