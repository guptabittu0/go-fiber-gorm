package auth

// LoginRequest represents the request for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents the request for registration
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RefreshTokenRequest represents the request for refreshing a token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ResetPasswordRequest represents the request for resetting a password
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ChangePasswordRequest represents the request for changing a password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// TokenResponse represents the response containing tokens
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // in seconds
	TokenType    string `json:"token_type"` // typically "Bearer"
}

// AuthResponse represents the authenticated user response
type AuthResponse struct {
	User  UserInfo      `json:"user"`
	Token TokenResponse `json:"token"`
}

// UserInfo represents user information in auth responses
type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
