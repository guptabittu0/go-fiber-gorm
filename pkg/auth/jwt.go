package auth

import (
	"fmt"
	"go-fiber-gorm/config"
	"go-fiber-gorm/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims represents custom JWT claims
type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	SecretKey string
	ExpiresIn time.Duration
}

// GenerateToken generates a new JWT token
func GenerateToken(userID uint, role string, cfg config.JWTConfig) (string, error) {
	// Parse the expiration duration
	expDuration, err := time.ParseDuration(cfg.ExpiresIn)
	if err != nil {
		return "", fmt.Errorf("invalid JWT expiration duration: %w", err)
	}

	// Set claims
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string, secretKey string) (*TokenClaims, error) {
	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewUnauthorizedError("invalid token signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid or expired token")
	}

	// Get claims from token
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.NewUnauthorizedError("invalid token claims")
	}

	return claims, nil
}

// GetUserIDFromToken extracts the user ID from JWT token
func GetUserIDFromToken(tokenString string, secretKey string) (uint, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetUserRoleFromToken extracts the user role from JWT token
func GetUserRoleFromToken(tokenString string, secretKey string) (string, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}
