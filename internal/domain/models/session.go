package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"primaryKey;autoIncrement"`
	JTI       uuid.UUID `gorm:"uniqueIndex;not null"`
	UserID    uuid.UUID `gorm:"not null;index"`
	IpAddress string    `gorm:"not null"`
	UserAgent string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
