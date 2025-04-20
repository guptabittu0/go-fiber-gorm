package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-fiber-gorm/config"
	"go-fiber-gorm/core/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client is the Redis client
var Client *redis.Client

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info("Connecting to Redis at", addr)

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Check if connection is successful
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis successfully")
	return Client, nil
}

// Get retrieves a value from cache by key
func Get[T any](key string) (T, bool) {
	var value T

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := Client.Get(ctx, key).Result()
	if err != nil {
		// Key does not exist or other error
		return value, false
	}

	if err := json.Unmarshal([]byte(data), &value); err != nil {
		logger.Error("Failed to unmarshal cached value:", err)
		return value, false
	}

	return value, true
}

// Set stores a value in cache with TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("Failed to marshal value for caching:", err)
		return err
	}

	return Client.Set(ctx, key, data, ttl).Err()
}

// Delete removes a value from cache
func Delete(keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return Client.Del(ctx, keys...).Err()
}

// Flush clears the entire cache
func Flush() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return Client.FlushAll(ctx).Err()
}

func Wrapper[T any](fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		data, err := fn()
		if err != nil {
			return data, err
		}

		return data, nil
	}
}
