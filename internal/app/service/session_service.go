package service

import (
	"context"
	"errors"
	"mytro-backend-auth/internal/app/repository"
	"mytro-backend-auth/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUserNotFoundOrInvalidPassword = errors.New("user not found or invalid password")
)

type SessionService struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
	tokenService      *TokensService // Assuming you have a TokensService for token generation
	logger            *zap.Logger
}

func NewSessionService(
	sessionRepository repository.SessionRepository,
	userRepository repository.UserRepository,
	tokenService *TokensService, // Assuming you have a TokensService for token generation
	logger *zap.Logger,
) *SessionService {
	return &SessionService{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		tokenService:      tokenService,
		logger:            logger,
	}
}

func (s *SessionService) SignInByEmail(ctx context.Context, usernameOrEmail, password string) (*models.Tokens, error) {
	user, err := s.userRepository.FindByEmailOrUsername(ctx, usernameOrEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrorUserNotFoundOrInvalidPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrorUserNotFoundOrInvalidPassword
		}
		return nil, err
	}

	tokens, err := s.tokenService.GenerateTokens(jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
	})
	if err != nil {
		return nil, err
	}

	s.sessionRepository.Create(ctx, &models.Session{
		UserID: user.ID,
		JTI:    tokens.JTI,
		// Todo: Add IP address and user agent
		ExpiresAt: tokens.RefreshTokenExpiresAt,
	})

	return tokens, nil
}
