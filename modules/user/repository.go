package user

import (
	"go-fiber-gorm/core/errors"

	"gorm.io/gorm"
)

// Repository handles database operations for users
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// Create creates a new user
func (r *Repository) Create(user *User) error {
	return r.DB.Create(user).Error
}

// FindByID finds a user by ID
func (r *Repository) FindByID(id uint) (*User, error) {
	var user User
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
func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
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
func (r *Repository) Update(user *User) error {
	return r.DB.Save(user).Error
}

// Delete deletes a user
func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&User{}, id).Error
}

// FindAll returns all users with pagination
func (r *Repository) FindAll(page, limit int) ([]User, int64, error) {
	var users []User
	var count int64

	// Count total records
	if err := r.DB.Model(&User{}).Count(&count).Error; err != nil {
		return nil, 0, errors.NewInternalServerError(err.Error())
	}

	// Get paginated records
	offset := (page - 1) * limit
	if err := r.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, errors.NewInternalServerError(err.Error())
	}

	return users, count, nil
}
