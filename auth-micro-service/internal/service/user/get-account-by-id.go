package user

import (
	"auth-micro-service/internal/dto"
	"context"

	"go.uber.org/zap"
)

func (s *Service) GetUserByAccountID(ctx context.Context, accountID string) (*dto.AccountAndUser, error) {
	userAccount, err := s.userRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		s.logger.Info("GetUserByAccountID", zap.Error(err))
		return nil, err
	}

	return userAccount, nil
}
