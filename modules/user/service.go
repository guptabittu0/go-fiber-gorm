package user

import (
	"go-fiber-gorm/core/errors"

	"github.com/go-playground/validator/v10"
)

// Service handles user-related business logic
type Service struct {
	repo      *Repository
	validator *validator.Validate
}

// NewService creates a new user service
func NewService(repo *Repository) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
	}
}

// Create creates a new user
func (s *Service) Create(req *CreateUserRequest) (*UserResponseDTO, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	// Check if user with this email already exists
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.NewBadRequestError("Email already in use")
	}

	// Create user
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Note: Password should be hashed in BeforeSave hook
		Role:     "user",       // Default role
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.NewInternalServerError("Failed to create user")
	}

	// Convert to DTO for response
	response := &UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// GetByID gets a user by ID
func (s *Service) GetByID(id uint) (*UserResponseDTO, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// Update updates a user
func (s *Service) Update(id uint, req *UpdateUserRequest) (*UserResponseDTO, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	// Get existing user
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already in use by another user
		existingUser, err := s.repo.FindByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.NewBadRequestError("Email already in use")
		}
		user.Email = req.Email
	}

	// Save updates
	if err := s.repo.Update(user); err != nil {
		return nil, errors.NewInternalServerError("Failed to update user")
	}

	response := &UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// Delete deletes a user
func (s *Service) Delete(id uint) error {
	// Check if user exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete user
	if err := s.repo.Delete(id); err != nil {
		return errors.NewInternalServerError("Failed to delete user")
	}

	return nil
}

// GetAll gets all users with pagination
func (s *Service) GetAll(page, limit int) ([]UserResponseDTO, int64, error) {
	// Default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	users, count, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response objects
	var responses []UserResponseDTO
	for _, user := range users {
		responses = append(responses, UserResponseDTO{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return responses, count, nil
}
