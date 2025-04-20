package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:100;not null" json:"name" validate:"required"`
	Email     string         `gorm:"size:100;not null;uniqueIndex" json:"email" validate:"required,email"`
	Password  string         `gorm:"size:100;not null" json:"-" validate:"required,min=6"`
	Role      string         `gorm:"size:20;not null;default:'user'" json:"role"`
}

// UserResponse is the response returned to clients
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeSave will be called before creating/updating a user
// Here you would hash the password
func (u *User) BeforeSave(tx *gorm.DB) error {
	// TODO: Implement password hashing
	return nil
}

// ToResponse converts a user to a response
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
