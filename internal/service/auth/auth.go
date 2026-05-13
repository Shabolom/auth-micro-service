package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/redis"
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthRepo interface {
	Register(ctx context.Context, reg *dto.Register) error
	CreateRefreshToken(ctx context.Context, session *dto.RefreshToken) error
	GetByEmail(ctx context.Context, email string) (*dto.Register, error)
	Logout(ctx context.Context, id uuid.UUID) error
	GetActiveRefreshToken(ctx context.Context, jti string, userID string, userAgent string) (string, error)
	UpdateRefreshTokenByID(ctx context.Context, oldJTI uuid.UUID, session *dto.RefreshToken) error
}

type Redis interface {
	RevokeSession(ctx context.Context, key string) error
	NewSession(userID string) *redis.Session
	CheckSessionStatus(ctx context.Context, jti string) error
	SaveSession(ctx context.Context, key string, value *redis.Session, expiration time.Duration) error
}

type RabbitMQ interface {
	Publish(ctx context.Context, routingKey string, contentType string, body []byte) error
}
type Service struct {
	authRepo AuthRepo

	rabbitMQ RabbitMQ

	redis Redis

	secret string
	logger *zap.Logger
}

func New(authRepo AuthRepo, rabbitMQ RabbitMQ, redis Redis, secret string, logger *zap.Logger) *Service {
	return &Service{
		authRepo: authRepo,
		rabbitMQ: rabbitMQ,
		secret:   secret,
		logger:   logger,
		redis:    redis,
	}
}
