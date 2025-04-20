package user

import "time"

// CreateUserRequest is the request to create a user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest is the request to update a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty,min=2"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponseDTO represents the user response for API
type UserResponseDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersResponseDTO represents a paginated list of users
type UsersResponseDTO struct {
	Users []UserResponseDTO `json:"users"`
	Meta  PaginationMeta    `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Pages int64 `json:"pages"`
}
