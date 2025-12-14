package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the domain model.
type User struct {
	ID            uuid.UUID
	Username      string
	Email         string
	Password      string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
