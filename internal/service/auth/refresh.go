package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) Refresh(ctx context.Context, oldRefToken string, userAgent string, ip string) (*dto.Tokens, error) {
	oldRefreshTokenClaims, err := utils.ParseToken(oldRefToken, s.secret, s.logger)
	if err != nil {
		s.logger.Warn("Parse refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	refreshTokenHash, err := s.authRepo.GetActiveRefreshToken(ctx, oldRefreshTokenClaims.ID, oldRefreshTokenClaims.UserID, userAgent)
	if err != nil {
		s.logger.Warn("Get active refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	err = utils.Compare(refreshTokenHash, oldRefToken)
	if err != nil {
		s.logger.Warn("Compare refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	oldJTIUUID, err := uuid.Parse(oldRefreshTokenClaims.ID)
	if err != nil {
		s.logger.Warn("Parse old jti", zap.Error(err))
		return &dto.Tokens{}, err
	}

	now := time.Now()

	newAccessJTI := uuid.New()
	newAccessToken, err := utils.GenerateAccessToken(oldRefreshTokenClaims.UserID, s.secret, newAccessJTI.String())
	if err != nil {
		s.logger.Warn("Generate access token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	newRefreshJTI := uuid.New()
	newRefreshToken, err := utils.GenerateRefreshToken(oldRefreshTokenClaims.UserID, s.secret, newRefreshJTI.String())
	if err != nil {
		s.logger.Warn("Generate refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	newRefreshTokenHash, err := utils.Hash(newRefreshToken)
	if err != nil {
		s.logger.Warn("Hash refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	userID, err := uuid.Parse(oldRefreshTokenClaims.UserID)
	if err != nil {
		s.logger.Warn("Parse old jti", zap.Error(err))
		return &dto.Tokens{}, err
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
		return &dto.Tokens{}, err
	}

	session := s.inmemorystorage.NewSession(userID.String())
	s.inmemorystorage.Save(newAccessJTI.String(), session)

	return &dto.Tokens{
		RefreshToken: newRefreshToken,
		AccessToken:  newAccessToken,
	}, nil
}
