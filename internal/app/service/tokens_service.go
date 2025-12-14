package service

import (
	"crypto/x509"
	"encoding/pem"
	"mytro-backend-auth/internal/domain/models"
	"mytro-backend-auth/internal/infrastructure/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokensService struct {
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration

	privateKey any
	publicKey  any

	logger *zap.Logger
}

func NewTokensService(keys *config.TokenConfig, logger *zap.Logger) *TokensService {
	privateKeyBlock, _ := pem.Decode(([]byte)(keys.PrivateKey()))
	publicKeyBlock, _ := pem.Decode(([]byte)(keys.PublicKey()))

	if privateKeyBlock == nil || publicKeyBlock == nil {
		logger.Fatal("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		logger.Fatal("failed to parse private key: " + err.Error())
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		logger.Fatal("failed to parse public key: " + err.Error())
	}

	return &TokensService{
		accessTokenLifetime:  keys.PublicTokenLifetime,
		refreshTokenLifetime: keys.PrivateTokenLifetime,

		privateKey: privateKey,
		publicKey:  publicKey,

		logger: logger,
	}
}

// GenerateTokens generates access and refresh tokens for a given set of claims
func (s *TokensService) GenerateTokens(claims jwt.Claims) (*models.Tokens, error) {
	now := time.Now()

	// Set claims for access token
	accessTokenClaims := claims.(jwt.MapClaims)
	accessTokenClaims["exp"] = now.Add(s.accessTokenLifetime).Unix()
	accessTokenClaims["iat"] = now.Unix()
	accessTokenClaims["jti"] = uuid.New().String()
	accessTokenClaims["type"] = "access"

	// Sign access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		s.logger.Error("failed to sign access token", zap.Error(err))
		return nil, err
	}

	expirationTime := now.Add(s.refreshTokenLifetime)

	// Set claims for refresh token
	refreshTokenClaims := jwt.MapClaims{
		"sub":  accessTokenClaims["sub"],
		"exp":  expirationTime.Unix(),
		"iat":  now.Unix(),
		"jti":  accessTokenClaims["jti"],
		"type": "refresh",
	}

	// Sign refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		s.logger.Error("failed to sign refresh token", zap.Error(err))
		return nil, err
	}

	return &models.Tokens{
		AccessToken:           accessTokenString,
		RefreshToken:          refreshTokenString,
		JTI:                   uuid.MustParse(accessTokenClaims["jti"].(string)),
		RefreshTokenExpiresAt: expirationTime,
	}, nil
}

func (s *TokensService) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			s.logger.Error("unexpected signing method", zap.String("method", token.Header["alg"].(string)))
			return nil, jwt.ErrTokenMalformed
		}
		return s.publicKey, nil
	})
	if err != nil {
		s.logger.Error("failed to parse token", zap.Error(err))
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "refresh" {
			s.logger.Error("invalid token type", zap.String("type", claims["type"].(string)))
			return uuid.Nil, jwt.ErrTokenMalformed
		}

		jtiStr, ok := claims["jti"].(string)
		if !ok {
			s.logger.Error("jti claim is missing or invalid")
			return uuid.Nil, jwt.ErrTokenMalformed
		}

		jti, err := uuid.Parse(jtiStr)
		if err != nil {
			s.logger.Error("failed to parse jti claim", zap.Error(err))
			return uuid.Nil, err
		}

		return jti, nil
	} else {
		s.logger.Error("invalid token claims")
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}
}
