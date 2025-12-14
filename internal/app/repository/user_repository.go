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

// UserRepository is an interface for interacting with the user repository.
type UserRepository interface {
	Save(ctx context.Context, user *models.User) error
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmailOrUsername(ctx context.Context, usernameOrEmail string) (*models.User, error)
}

// userRepository is a concrete implementation of the UserRepository interface.
type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// Save saves a user to the repository.
// If the user is nil, it returns an error.
func (r *userRepository) Save(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	userDTO := r.toDTO(*user)

	return r.db.WithContext(ctx).Save(&userDTO).Error
}

// Create creates a new user by email.
// If the user is nil, it returns an error.
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	userDTO := r.toDTO(*user)

	return r.db.WithContext(ctx).Create(&userDTO).Error
}

// FindByID finds a user by ID.
// If the user is not found, it returns nil.
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user dto.UserDTO
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	userDomain := r.toDomain(user)
	return &userDomain, nil
}

func (r *userRepository) FindByEmailOrUsername(ctx context.Context, usernameOrEmail string) (*models.User, error) {
	var user dto.UserDTO
	if err := r.db.WithContext(ctx).First(&user, "email = ? OR username = ?", usernameOrEmail, usernameOrEmail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	userDomain := r.toDomain(user)
	return &userDomain, nil
}

// converts UserDTO to User domain model.
func (r *userRepository) toDomain(dto dto.UserDTO) models.User {
	return models.User{
		ID:            dto.ID,
		Username:      dto.Username,
		Email:         dto.Email,
		Password:      dto.PasswordHash,
		EmailVerified: dto.EmailVerified,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

// converts User domain model to UserDTO.
func (r *userRepository) toDTO(user models.User) dto.UserDTO {
	return dto.UserDTO{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		PasswordHash:  user.Password,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}
