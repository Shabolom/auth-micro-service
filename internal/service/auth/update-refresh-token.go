package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/shortcut"
	"auth-micro-service/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) UpdateRefreshToken(ctx context.Context, oldTokens *dto.Tokens, userAgent string, ip string) (dto.Tokens, error) {
	oldRefreshTokenClaims, err := utils.ParseToken(oldTokens.RefreshToken, s.secret, s.logger)
	if err != nil {
		s.logger.Info("Parse refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	oldAccessTokenClaims, err := utils.ParseToken(oldTokens.AccessToken, s.secret, s.logger)
	if err != nil {
		s.logger.Info("Parse access token", zap.Error(err))
		return dto.Tokens{}, err
	}

	if oldAccessTokenClaims.UserID != oldRefreshTokenClaims.UserID {
		s.logger.Info("user ID mismatch")
		return dto.Tokens{}, errors.New("access token and refresh token belong to different users")
	}

	userID, err := uuid.Parse(oldRefreshTokenClaims.UserID)
	if err != nil {
		s.logger.Info("Parse user id", zap.Error(err))
		return dto.Tokens{}, err
	}

	refreshTokenHash, err := s.authRepo.GetActiveRefreshToken(ctx, oldRefreshTokenClaims.ID, oldRefreshTokenClaims.UserID, userAgent)
	if err != nil {
		if !errors.Is(err, shortcut.ErrNoRows) {
			s.logger.Info("Get active refresh token", zap.Error(err))
			return dto.Tokens{}, err
		}
	}

	if refreshTokenHash != "" {
		err = utils.Compare(refreshTokenHash, oldTokens.RefreshToken)
		if err != nil {
			s.logger.Info("Compare refresh token", zap.Error(err))
			return dto.Tokens{}, err
		}
	}

	oldJTIUUID, err := uuid.Parse(oldRefreshTokenClaims.ID)
	if err != nil {
		s.logger.Info("Parse old jti", zap.Error(err))
		return dto.Tokens{}, err
	}

	now := time.Now()

	newAccessJTI := uuid.New()
	accessToken, err := utils.GenerateAccessToken(userID.String(), s.secret, newAccessJTI.String())
	if err != nil {
		s.logger.Info("Generate access token", zap.Error(err))
		return dto.Tokens{}, err
	}

	newRefreshJTI := uuid.New()
	refreshToken, err := utils.GenerateRefreshToken(userID.String(), s.secret, newRefreshJTI.String())
	if err != nil {
		s.logger.Info("Generate refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	newRefreshTokenHash, err := utils.Hash(refreshToken)
	if err != nil {
		s.logger.Info("Hash refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	storageRefreshToken := &dto.RefreshToken{
		ID:        newRefreshJTI,
		UserID:    userID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: now.Add(72 * time.Hour),
		RevokedAt: nil,
		CreatedAt: now,
		UserAgent: userAgent,
		IP:        ip,
	}

	err = s.authRepo.UpdateRefreshTokenByID(ctx, oldJTIUUID, storageRefreshToken)
	if err != nil {
		s.logger.Info("Update refresh token by id", zap.Error(err))
		return dto.Tokens{}, err
	}

	err = s.redis.RevokeSession(ctx, oldRefreshTokenClaims.ID)
	if err != nil {
		s.logger.Info("Revoke old refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	session := s.redis.NewSession(userID.String())
	err = s.redis.SaveSession(ctx, newAccessJTI.String(), session, time.Minute*15)
	if err != nil {
		s.logger.Info("Save new access jti", zap.Error(err))
		return dto.Tokens{}, err
	}

	return dto.Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
