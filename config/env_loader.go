package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Environment types
const (
	Development = "development"
	Testing     = "testing"
	Production  = "production"
)

// LoadEnvForCurrentEnvironment loads environment variables based on the current environment
func LoadEnvForCurrentEnvironment() error {
	// Get environment from ENV variable, default to development
	env := os.Getenv("ENV")
	if env == "" {
		env = Development
	}

	// List of .env files to try loading, in order
	var envFiles []string

	// Always try to load base .env file
	envFiles = append(envFiles, ".env")

	// Load environment-specific file if it exists
	if env != Development {
		envFiles = append(envFiles, fmt.Sprintf(".env.%s", env))
	}

	// Load environment-specific local overrides if they exist
	localEnvFile := fmt.Sprintf(".env.%s.local", env)
	if fileExists(localEnvFile) {
		envFiles = append(envFiles, localEnvFile)
	}

	// Local overrides for all environments
	if fileExists(".env.local") {
		envFiles = append(envFiles, ".env.local")
	}

	// Try loading each file in sequence
	for _, file := range envFiles {
		if fileExists(file) {
			if err := godotenv.Load(file); err != nil {
				return fmt.Errorf("error loading %s: %w", file, err)
			}
		}
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateEnvFile creates a .env file from the example if it doesn't exist
func CreateEnvFile() error {
	// Check if .env file already exists
	if fileExists(".env") {
		return nil
	}

	// Check if .env.example exists
	if !fileExists(".env.example") {
		return fmt.Errorf(".env.example file not found")
	}

	// Read the content of .env.example
	content, err := os.ReadFile(".env.example")
	if err != nil {
		return fmt.Errorf("failed to read .env.example: %w", err)
	}

	// Create .env file with the same content
	if err := os.WriteFile(".env", content, 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	return nil
}

// GetConfigDir returns the absolute path to the configuration directory
func GetConfigDir() (string, error) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if we're in the project root or a subdirectory
	configDir := filepath.Join(wd, "config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// Try going up one directory (for when run from subdirectories)
		parent := filepath.Dir(wd)
		configDir = filepath.Join(parent, "config")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			return "", fmt.Errorf("config directory not found")
		}
	}

	return configDir, nil
}
