package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionDTO struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	JTI       uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	IpAddress string    `gorm:"type:varchar(255);not null"`
	UserAgent string    `gorm:"type:text;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}
