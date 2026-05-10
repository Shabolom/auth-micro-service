package user

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"context"

	"go.uber.org/zap"
)

type UserRepo interface {
	CreateUserDescription(ctx context.Context, accountID string, name *string, age *int) (*dto.UserDescription, error)
	CheckAccountByID(ctx context.Context, userID string) error
	GetAccountByID(ctx context.Context, accountID string) (*dto.AccountAndUser, error)
	GetAccounts(ctx context.Context) ([]dto.AccountAndUser, error)
	DeleteUser(ctx context.Context, accountID string) error
	UpdateUser(ctx context.Context, accountID string, req *authv1.UpdateUser) (*dto.AccountAndUser, error)
}

type Service struct {
	userRepo UserRepo
	secret   string
	logger   *zap.Logger
}

func New(userRepo UserRepo, secret string, logger *zap.Logger) *Service {
	return &Service{
		userRepo: userRepo,
		secret:   secret,
		logger:   logger,
	}
}
