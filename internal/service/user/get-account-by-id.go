package user

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"

	"go.uber.org/zap"
)

func (s *Service) GetUserByID(ctx context.Context, strAccessToken string) (*dto.AccountAndUser, error) {
	accessTokenClaims, err := utils.ParseToken(strAccessToken, s.secret, s.logger)
	if err != nil {
		s.logger.Info("Failed to parse token", zap.Error(err))
		return nil, err
	}

	userAccount, err := s.userRepo.GetUserByID(ctx, accessTokenClaims.UserID)
	if err != nil {
		s.logger.Info("GetUserByID", zap.Error(err))
		return nil, err
	}

	return userAccount, nil
}
