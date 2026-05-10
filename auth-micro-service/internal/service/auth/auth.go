package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/inmemory"
	"context"

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

type SessionStorage interface {
	NewSession(userID string) inmemory.Session
	Save(jti string, session inmemory.Session)
	Get(jti string) (inmemory.Session, bool)
	Revoke(jti string)
}
type Service struct {
	authRepo        AuthRepo
	inmemorystorage SessionStorage
	secret          string
	logger          *zap.Logger
}

func New(authRepo AuthRepo, inmemorystorage SessionStorage, secret string, logger *zap.Logger) *Service {
	return &Service{
		authRepo:        authRepo,
		inmemorystorage: inmemorystorage,
		secret:          secret,
		logger:          logger,
	}
}
