package repository

import (
	"context"
	"errors"
	"mytro-backend-auth/internal/domain/models"
	"mytro-backend-auth/internal/infrastructure/database/dto"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SessionRepository interface {
	FindByJTI(ctx context.Context, id uuid.UUID) (*models.Session, error)
	Create(ctx context.Context, session *models.Session) error
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
}

type sessionRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewSessionRepository(db *gorm.DB, logger *zap.Logger) SessionRepository {
	return &sessionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *sessionRepository) FindByJTI(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	if err := r.db.Where("jti = ?", id).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	sessionDTO := toSessionDTO(*session)

	return r.db.WithContext(ctx).Create(&sessionDTO).Error
}

func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	return r.db.WithContext(ctx).Save(session).Error
}

func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("jti = ?", id).Delete(&models.Session{}).Error
}

func (r *sessionRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

func toSessionDTO(session models.Session) dto.SessionDTO {
	return dto.SessionDTO{
		ID:        session.ID,
		JTI:       session.JTI,
		UserID:    session.UserID,
		IpAddress: session.IpAddress,
		UserAgent: session.UserAgent,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}
