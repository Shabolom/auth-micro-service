package user

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func (s *Service) DeleteUser(ctx context.Context, tokens *dto.Tokens) error {
	accessTokenClaims, err := utils.ParseToken(tokens.AccessToken, s.secret, s.logger)
	if err != nil {
		log.Infof("Error parsing access token claims: %s", err.Error())
		return err
	}

	refreshTokenClaims, err := utils.ParseToken(tokens.RefreshToken, s.secret, s.logger)
	if err != nil {
		log.Infof("Error parsing refresh token claims: %s", err.Error())
		return err
	}

	err = s.userRepo.DeleteUser(ctx, accessTokenClaims.UserID)
	fmt.Println(accessTokenClaims.UserID, 123123123)
	if err != nil {
		s.logger.Info("delete error")
		return fmt.Errorf("DeleteUser: %w", err)
	}

	jtiRefreshToken, err := uuid.Parse(refreshTokenClaims.ID)
	if err != nil {
		s.logger.Info("error parsing refresh token claims")
		return err
	}

	fmt.Println(jtiRefreshToken, 2222222222)
	err = s.refreshTokenRepo.Logout(ctx, jtiRefreshToken)
	if err != nil {
		s.logger.Info("logout error", zap.Error(err))
		return fmt.Errorf("Logout: %w", err)
	}

	s.inMemoryCache.Revoke(accessTokenClaims.UserID)

	return nil
}
