package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserDTO represents the user data transfer object for database operations.
type UserDTO struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username      string    `gorm:"type:varchar(32);uniqueIndex;not null"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string    `gorm:"type:text;not null"`
	EmailVerified bool      `gorm:"not null;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (UserDTO) TableName() string {
	return "users"
}
