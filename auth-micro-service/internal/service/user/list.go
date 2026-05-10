package user

import (
	"auth-micro-service/internal/dto"
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *Service) GetAccounts(ctx context.Context) ([]dto.AccountAndUser, error) {
	accounts, err := s.userRepo.GetAccounts(ctx)
	if err != nil {
		s.logger.Info("GetAccounts error", zap.Error(err))
		return nil, fmt.Errorf("GetAccounts: %w", err)
	}

	return accounts, nil
}
