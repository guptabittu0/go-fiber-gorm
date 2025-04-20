package auth

import (
	"go-fiber-gorm/core/errors"

	"gorm.io/gorm"
)

// Repository handles database operations for auth
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// CreateSession creates a new user session
func (r *Repository) CreateSession(session *Session) error {
	return r.DB.Create(session).Error
}

// FindSessionByToken finds a session by refresh token
func (r *Repository) FindSessionByToken(refreshToken string) (*Session, error) {
	var session Session
	err := r.DB.Where("refresh_token = ? AND is_blocked = ?", refreshToken, false).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Session")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &session, nil
}

// InvalidateSession marks a session as blocked
func (r *Repository) InvalidateSession(sessionID uint) error {
	return r.DB.Model(&Session{}).Where("id = ?", sessionID).Update("is_blocked", true).Error
}

// InvalidateAllUserSessions marks all sessions for a user as blocked
func (r *Repository) InvalidateAllUserSessions(userID uint) error {
	return r.DB.Model(&Session{}).Where("user_id = ?", userID).Update("is_blocked", true).Error
}

// DeleteExpiredSessions deletes all expired sessions
func (r *Repository) DeleteExpiredSessions() error {
	return r.DB.Where("expires_at < NOW()").Delete(&Session{}).Error
}
