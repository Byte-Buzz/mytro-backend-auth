package models

import (
	"time"

	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string

	JTI                   uuid.UUID
	RefreshTokenExpiresAt time.Time
}
