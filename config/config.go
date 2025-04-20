package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config stores all configuration of the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

// ServerConfig stores server related configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig stores database configuration
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
	SSLMode  string
}

// JWTConfig stores JWT configuration
type JWTConfig struct {
	Secret          string
	AccessExpiryIn  uint
	RefreshExpiryIn uint
}

// RedisConfig stores Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// LoadConfig reads configuration from .env file
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	dbPort, err := parseEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	redisPort, err := parseEnvInt("REDIS_PORT", 6379)
	if err != nil {
		return nil, err
	}

	redisDB, err := parseEnvInt("REDIS_DB", 0)
	if err != nil {
		return nil, err
	}

	accessExpiryIn, err := parseEnvUint("JWT_ACCESS_EXPIRY", 3600) // 1 hour
	if err != nil {
		return nil, err
	}

	refreshExpiryIn, err := parseEnvUint("JWT_REFRESH_EXPIRY", 604800) // 7 days
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "fiber_gorm_db"),
			Port:     dbPort,
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your_secret_key"),
			AccessExpiryIn:  uint(accessExpiryIn),
			RefreshExpiryIn: uint(refreshExpiryIn),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     redisPort,
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
	}, nil
}

// getEnv reads environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// parseEnvInt parses an integer environment variable with a default value
func parseEnvInt(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", key, err)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

// parseEnvUint parses an unsigned integer environment variable with a default value
func parseEnvUint(key string, defaultValue uint) (uint, error) {
	if value, exists := os.LookupEnv(key); exists {
		uintValue, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", key, err)
		}
		return uint(uintValue), nil
	}
	return defaultValue, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode)
}
