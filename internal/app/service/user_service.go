package service

import (
	"context"
	"mytro-backend-auth/internal/app/repository"
	"mytro-backend-auth/internal/domain/models"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService provides user-related services.
type UserService struct {
	userRepository repository.UserRepository
	logger         *zap.Logger
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepository repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (s *UserService) SignUpByEmail(ctx context.Context, user *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	user.Email = strings.ToLower(user.Email)
	user.Username = strings.ToLower(user.Username)

	return s.userRepository.Create(ctx, user)
}
