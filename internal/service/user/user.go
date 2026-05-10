package user

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, accountID string) (*dto.AccountAndUser, error)
	GetUsers(ctx context.Context) ([]*dto.AccountAndUser, error)
	DeleteUser(ctx context.Context, accountID string) error
	UpdateUser(ctx context.Context, accountID string, req *authv1.UpdateUser) (*dto.AccountAndUser, error)
}

type RefreshTokenRepo interface {
	Logout(ctx context.Context, id uuid.UUID) error
}

type InMemoryCache interface {
	Revoke(jti string)
}
type Service struct {
	userRepo         UserRepo
	refreshTokenRepo RefreshTokenRepo
	inMemoryCache    InMemoryCache
	secret           string
	logger           *zap.Logger
}

func New(refreshTokenRepo RefreshTokenRepo, userRepo UserRepo, inMemoryCache InMemoryCache, secret string, logger *zap.Logger) *Service {
	return &Service{
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
		secret:           secret,
		logger:           logger,
		inMemoryCache:    inMemoryCache,
	}
}
