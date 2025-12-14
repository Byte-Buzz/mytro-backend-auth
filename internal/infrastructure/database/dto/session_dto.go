package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionDTO struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	JTI       uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	IpAddress string    `gorm:"type:varchar(255);not null"`
	UserAgent string    `gorm:"type:text;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SessionDTO) TableName() string {
	return "sessions"
}
