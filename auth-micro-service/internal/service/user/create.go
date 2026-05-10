package user

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"

	"go.uber.org/zap"
)

func (s *Service) Create(ctx context.Context, request *dto.CreateUserDescriptionRequest, accessToken string) (*dto.UserDescription, error) {
	accessTokenClaims, err := utils.ParseToken(accessToken, s.secret, s.logger)
	if err != nil {
		s.logger.Info("Error parsing token", zap.Error(err))
		return &dto.UserDescription{}, err
	}

	err = s.userRepo.CheckAccountByID(ctx, accessTokenClaims.UserID)
	if err != nil {
		s.logger.Info("Error checking account", zap.Error(err))
		return &dto.UserDescription{}, err
	}

	user, err := s.userRepo.CreateUserDescription(ctx, accessTokenClaims.UserID, &request.Name, &request.Age)
	if err != nil {
		s.logger.Info("Error checking user", zap.Error(err))
		return &dto.UserDescription{}, err
	}

	return user, nil
}
