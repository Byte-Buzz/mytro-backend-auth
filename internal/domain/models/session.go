package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	JTI       uuid.UUID
	UserID    uuid.UUID
	IpAddress string
	UserAgent string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
