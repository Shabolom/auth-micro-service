package user

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/shortcut"
	"auth-micro-service/pkg/utils"
	"context"

	"go.uber.org/zap"
)

func (s *Service) UpdateUsers(ctx context.Context, strAccessToken string, req *authv1.UpdateUser) (*dto.AccountAndUser, error) {
	accessTokenClaims, err := utils.ParseToken(strAccessToken, s.secret, s.logger)
	if err != nil {
		s.logger.Info("Error parsing token", zap.Error(err))
		return &dto.AccountAndUser{}, err
	}

	err = updateUserValidate(req)
	if err != nil {
		s.logger.Info("Error updating user", zap.Error(err))
		return &dto.AccountAndUser{}, err
	}

	userAndAccount, err := s.userRepo.UpdateUser(ctx, accessTokenClaims.UserID, req)
	if err != nil {
		s.logger.Info("Error updating user", zap.Error(err))
		return &dto.AccountAndUser{}, err
	}

	return userAndAccount, nil
}

func updateUserValidate(req *authv1.UpdateUser) error {

	if req.GetAge() == 0 {
		return shortcut.ErrEmptyFields
	} else if req.GetAge() < 13 {
		return shortcut.ErrAgeLimit
	}
	if req.GetName() == "" {
		return shortcut.ErrEmptyFields
	}
	if req.GetMail() == "" {
		return shortcut.ErrEmptyFields
	}

	return nil
}
