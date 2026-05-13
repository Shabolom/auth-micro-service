package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) Logout(ctx context.Context, userAgent string, tokens *dto.Tokens) error {
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

	refreshTokenHash, err := s.authRepo.GetActiveRefreshToken(ctx, refreshTokenClaims.ID, refreshTokenClaims.UserID, userAgent)
	if err != nil {
		s.logger.Warn("Couldn't get active refresh token")
		return err
	}

	err = utils.Compare(refreshTokenHash, tokens.RefreshToken)
	if err != nil {
		s.logger.Warn("Couldn't compare refresh token")
		return err
	}

	if err = s.authRepo.Logout(ctx, refreshTokenID); err != nil {
		s.logger.Info(err.Error())
		return err
	}

	err = s.redis.RevokeSession(ctx, accessTokenClaims.ID)
	if err != nil {
		s.logger.Info("revoke err:", zap.Error(err))
		return err
	}

	return nil
}
