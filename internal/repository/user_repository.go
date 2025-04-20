package repository

import (
	"go-fiber-gorm/internal/model"
	"go-fiber-gorm/pkg/errors"

	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("User")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("User")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&model.User{}, id).Error
}

// FindAll returns all users with pagination
func (r *UserRepository) FindAll(page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var count int64

	// Count total records
	if err := r.DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return nil, 0, errors.NewInternalServerError(err.Error())
	}

	// Get paginated records
	offset := (page - 1) * limit
	if err := r.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, errors.NewInternalServerError(err.Error())
	}

	return users, count, nil
}
