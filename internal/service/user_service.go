package service

import (
	"go-fiber-gorm/internal/model"
	"go-fiber-gorm/internal/repository"
	"go-fiber-gorm/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo  *repository.UserRepository
	validator *validator.Validate
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// Create creates a new user
func (s *UserService) Create(req *model.CreateUserRequest) (*model.UserResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	// Check if user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.NewBadRequestError("Email already in use")
	}

	// Create user
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Note: Password should be hashed in BeforeSave hook
		Role:     "user",       // Default role
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.NewInternalServerError("Failed to create user")
	}

	return user.ToResponse(), nil
}

// GetByID gets a user by ID
func (s *UserService) GetByID(id uint) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

// Update updates a user
func (s *UserService) Update(id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	// Get existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already in use by another user
		existingUser, err := s.userRepo.FindByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.NewBadRequestError("Email already in use")
		}
		user.Email = req.Email
	}

	// Save updates
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.NewInternalServerError("Failed to update user")
	}

	return user.ToResponse(), nil
}

// Delete deletes a user
func (s *UserService) Delete(id uint) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete user
	if err := s.userRepo.Delete(id); err != nil {
		return errors.NewInternalServerError("Failed to delete user")
	}

	return nil
}

// GetAll gets all users with pagination
func (s *UserService) GetAll(page, limit int) ([]model.UserResponse, int64, error) {
	// Default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	users, count, err := s.userRepo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response objects
	var responses []model.UserResponse
	for _, user := range users {
		responses = append(responses, *user.ToResponse())
	}

	return responses, count, nil
}
