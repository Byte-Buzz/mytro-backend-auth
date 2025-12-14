package service

import (
	"context"
	"errors"
	"mytro-backend-auth/internal/app/repository"
	"mytro-backend-auth/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (s *SessionService) SignInByEmail(ctx context.Context, usernameOrEmail, password, ip, userAgent string) (*models.Tokens, error) {
	user, err := s.userRepository.FindByEmailOrUsername(ctx, usernameOrEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrorUserNotFoundOrInvalidPassword
	}

	s.logger.Debug("User found", zap.String("userID", user.ID.String()))
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrorUserNotFoundOrInvalidPassword
		}
		return nil, err
	}

	s.logger.Debug("Password verified", zap.String("userID", user.ID.String()))
	tokens, err := s.tokenService.GenerateTokens(jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Tokens generated", zap.String("userID", user.ID.String()), zap.String("jti", tokens.JTI.String()))
	s.sessionRepository.Create(ctx, &models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		JTI:       tokens.JTI,
		IpAddress: ip,
		UserAgent: userAgent,
		ExpiresAt: tokens.RefreshTokenExpiresAt,
	})

	return tokens, nil
}

func (s *SessionService) SignOut(ctx context.Context, token string) error {
	jti, err := s.tokenService.ValidateRefreshToken(token)
	if err != nil {
		return err
	}

	return s.sessionRepository.DeleteByJTI(ctx, jti)
}

func (s *SessionService) RefreshSession(ctx context.Context, token, ip, userAgent string) (*models.Tokens, error) {
	jti, err := s.tokenService.ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionRepository.FindByJTI(ctx, jti)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("session not found")
	}

	user, err := s.userRepository.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	tokens, err := s.tokenService.GenerateTokens(jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
	})
	if err != nil {
		return nil, err
	}

	session.JTI = tokens.JTI
	session.IpAddress = ip
	session.UserAgent = userAgent
	session.ExpiresAt = tokens.RefreshTokenExpiresAt

	err = s.sessionRepository.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
