package auth

import (
	"time"

	"go-fiber-gorm/modules/user"
)

// Claims represents the JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// Session represents a user session
type Session struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	RefreshToken string    `gorm:"size:255;not null;uniqueIndex" json:"-"`
	UserAgent    string    `gorm:"size:255;not null" json:"user_agent"`
	ClientIP     string    `gorm:"size:100;not null" json:"client_ip"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	IsBlocked    bool      `gorm:"default:false;not null" json:"is_blocked"`
}

// TokenDetails contains both access and refresh tokens
type TokenDetails struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessUUID   string    `json:"-"`
	RefreshUUID  string    `json:"-"`
	AtExpires    int64     `json:"-"`
	RtExpires    int64     `json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// AuthenticatedUser represents a user who has been authenticated
type AuthenticatedUser struct {
	User  *user.User    `json:"-"`
	Token *TokenDetails `json:"-"`
}
