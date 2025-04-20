package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go-fiber-gorm/core/errors"
	"go-fiber-gorm/modules/user"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Service handles auth-related business logic
type Service struct {
	repo          *Repository
	userRepo      *user.Repository
	validator     *validator.Validate
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// ServiceConfig contains configuration for the auth service
type ServiceConfig struct {
	JWTSecret     string
	AccessExpiry  time.Duration // Usually short, e.g., 15 minutes
	RefreshExpiry time.Duration // Usually longer, e.g., 7 days
}

// NewService creates a new auth service
func NewService(repo *Repository, userRepo *user.Repository, config ServiceConfig) *Service {
	return &Service{
		repo:          repo,
		userRepo:      userRepo,
		validator:     validator.New(),
		jwtSecret:     config.JWTSecret,
		accessExpiry:  config.AccessExpiry,
		refreshExpiry: config.RefreshExpiry,
	}
}

// Register registers a new user
func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	// Check if user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.NewBadRequestError("Email already in use")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to hash password")
	}

	// Create user
	newUser := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user", // Default role
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, errors.NewInternalServerError("Failed to create user")
	}

	// Generate tokens
	tokenDetails, err := s.generateTokens(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate tokens")
	}

	// Save session
	session := &Session{
		UserID:       newUser.ID,
		RefreshToken: tokenDetails.RefreshToken,
		UserAgent:    "Not provided", // Should be extracted from request context in actual implementation
		ClientIP:     "Not provided", // Should be extracted from request context in actual implementation
		ExpiresAt:    time.Unix(tokenDetails.RtExpires, 0),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, errors.NewInternalServerError("Failed to create session")
	}

	// Prepare response
	response := &AuthResponse{
		User: UserInfo{
			ID:    newUser.ID,
			Name:  newUser.Name,
			Email: newUser.Email,
			Role:  newUser.Role,
		},
		Token: TokenResponse{
			AccessToken:  tokenDetails.AccessToken,
			RefreshToken: tokenDetails.RefreshToken,
			ExpiresIn:    tokenDetails.AtExpires - time.Now().Unix(),
			TokenType:    "Bearer",
		},
	}

	return response, nil
}

// Login authenticates a user
func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	// Find user by email
	foundUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password)); err != nil {
		return nil, errors.NewUnauthorizedError("Invalid credentials")
	}

	// Generate tokens
	tokenDetails, err := s.generateTokens(foundUser.ID, foundUser.Email, foundUser.Role)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate tokens")
	}

	// Save session
	session := &Session{
		UserID:       foundUser.ID,
		RefreshToken: tokenDetails.RefreshToken,
		UserAgent:    "Not provided", // Should be extracted from request context
		ClientIP:     "Not provided", // Should be extracted from request context
		ExpiresAt:    time.Unix(tokenDetails.RtExpires, 0),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, errors.NewInternalServerError("Failed to create session")
	}

	// Prepare response
	response := &AuthResponse{
		User: UserInfo{
			ID:    foundUser.ID,
			Name:  foundUser.Name,
			Email: foundUser.Email,
			Role:  foundUser.Role,
		},
		Token: TokenResponse{
			AccessToken:  tokenDetails.AccessToken,
			RefreshToken: tokenDetails.RefreshToken,
			ExpiresIn:    tokenDetails.AtExpires - time.Now().Unix(),
			TokenType:    "Bearer",
		},
	}

	return response, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *Service) RefreshToken(req *RefreshTokenRequest) (*TokenResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	// Find session by refresh token
	session, err := s.repo.FindSessionByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid refresh token")
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		// Invalidate session
		_ = s.repo.InvalidateSession(session.ID)
		return nil, errors.NewUnauthorizedError("Refresh token expired")
	}

	// Find user associated with the session
	foundUser, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to find user")
	}

	// Generate new tokens
	tokenDetails, err := s.generateTokens(foundUser.ID, foundUser.Email, foundUser.Role)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate tokens")
	}

	// Invalidate old session
	if err := s.repo.InvalidateSession(session.ID); err != nil {
		return nil, errors.NewInternalServerError("Failed to invalidate old session")
	}

	// Create new session
	newSession := &Session{
		UserID:       foundUser.ID,
		RefreshToken: tokenDetails.RefreshToken,
		UserAgent:    session.UserAgent, // Preserve user agent
		ClientIP:     session.ClientIP,  // Preserve client IP
		ExpiresAt:    time.Unix(tokenDetails.RtExpires, 0),
	}

	if err := s.repo.CreateSession(newSession); err != nil {
		return nil, errors.NewInternalServerError("Failed to create new session")
	}

	// Return new tokens
	return &TokenResponse{
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
		ExpiresIn:    tokenDetails.AtExpires - time.Now().Unix(),
		TokenType:    "Bearer",
	}, nil
}

// Logout invalidates a session
func (s *Service) Logout(refreshToken string) error {
	// Find session by refresh token
	session, err := s.repo.FindSessionByToken(refreshToken)
	if err != nil {
		// Token might be already invalid, so don't return error
		return nil
	}

	// Invalidate session
	return s.repo.InvalidateSession(session.ID)
}

// LogoutAll invalidates all sessions for a user
func (s *Service) LogoutAll(userID uint) error {
	return s.repo.InvalidateAllUserSessions(userID)
}

// ChangePassword changes a user's password
func (s *Service) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return errors.NewValidationError(err)
	}

	// Find user
	foundUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.OldPassword)); err != nil {
		return errors.NewBadRequestError("Current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalServerError("Failed to hash password")
	}

	// Update password
	foundUser.Password = string(hashedPassword)
	if err := s.userRepo.Update(foundUser); err != nil {
		return errors.NewInternalServerError("Failed to update password")
	}

	// Invalidate all sessions for security
	return s.repo.InvalidateAllUserSessions(userID)
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.NewUnauthorizedError(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if token is expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.NewUnauthorizedError("Token expired")
		}

		// Extract claims
		userClaims := &Claims{
			UserID: uint(claims["user_id"].(float64)),
			Email:  claims["email"].(string),
			Role:   claims["role"].(string),
		}

		return userClaims, nil
	}

	return nil, errors.NewUnauthorizedError("Invalid token")
}

// generateTokens generates access and refresh tokens
func (s *Service) generateTokens(userID uint, email, role string) (*TokenDetails, error) {
	now := time.Now()

	td := &TokenDetails{
		AtExpires:   now.Add(s.accessExpiry).Unix(),
		RtExpires:   now.Add(s.refreshExpiry).Unix(),
		AccessUUID:  generateUUID(),
		RefreshUUID: generateUUID(),
		ExpiresAt:   now.Add(s.accessExpiry),
	}

	// Create access token
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"uuid":    td.AccessUUID,
		"exp":     td.AtExpires,
		"iat":     now.Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	var err error
	td.AccessToken, err = accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"uuid":    td.RefreshUUID,
		"exp":     td.RtExpires,
		"iat":     now.Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// generateUUID generates a random UUID
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
