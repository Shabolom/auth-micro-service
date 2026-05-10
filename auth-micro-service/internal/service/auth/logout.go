package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *Service) Logout(ctx context.Context, tokens *dto.Tokens) error {
	refreshTokenClaims, err := utils.ParseToken(tokens.RefreshToken, s.secret, s.logger)
	if err != nil {
		s.logger.Info(err.Error())
		return err
	}

	accessTokenClaims, err := utils.ParseToken(tokens.AccessToken, s.secret, s.logger)
	if err != nil {
		s.logger.Info(err.Error())
		return err
	}

	if accessTokenClaims.UserID != refreshTokenClaims.UserID {
		s.logger.Info("user ID mismatch")
		return errors.New("access token and refresh token belong to different users")
	}

	refreshTokenID, err := uuid.Parse(refreshTokenClaims.ID)
	if err != nil {
		s.logger.Info("Couldn't parse refresh token id")
		return err
	}

	if err = s.authRepo.Logout(ctx, refreshTokenID); err != nil {
		s.logger.Info(err.Error())
		return err
	}

	s.inmemorystorage.Revoke(accessTokenClaims.ID)

	return nil
}
