package user

import (
	"auth-micro-service/internal/dto"
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *Service) GetUsers(ctx context.Context) ([]*dto.AccountAndUser, error) {
	accounts, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		s.logger.Info("GetUsers error", zap.Error(err))
		return nil, fmt.Errorf("GetUsers: %w", err)
	}

	return accounts, nil
}
