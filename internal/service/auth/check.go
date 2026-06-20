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

func (s *Service) Check(ctx context.Context, tokens *dto.Tokens, userAgent string, ip string) (*dto.Tokens, error) {
	s.logger.Info("check started",
		zap.String("user_agent", userAgent),
		zap.String("ip", ip),
	)

	refTokenClaims, err := utils.ParseTokenForCheck(tokens.RefreshToken, s.secret, s.logger)
	if err != nil {
		s.logger.Warn("refresh token parse failed",
			zap.Error(err),
		)

		if !errors.Is(err, shortcut.ErrInvalidToken) {
			s.logger.Warn("refresh token is invalid")
			return nil, shortcut.ErrInvalidToken
		}
	}

	s.logger.Info("refresh token parsed",
		zap.String("refresh_jti", refTokenClaims.ID),
		zap.String("user_id", refTokenClaims.UserID),
	)

	accessTokenClaims, err := utils.ParseTokenForCheck(tokens.AccessToken, s.secret, s.logger)
	if err != nil {
		s.logger.Warn("access token parse failed",
			zap.Error(err),
		)

		if !errors.Is(err, shortcut.ErrInvalidToken) {
			s.logger.Warn("access token is invalid")
			return nil, shortcut.ErrInvalidToken
		}
	}

	s.logger.Info("access token parsed",
		zap.String("access_jti", accessTokenClaims.ID),
		zap.String("user_id", accessTokenClaims.UserID),
	)

	if refTokenClaims.UserID != accessTokenClaims.UserID {
		s.logger.Warn("token pair belongs to different users",
			zap.String("refresh_user_id", refTokenClaims.UserID),
			zap.String("access_user_id", accessTokenClaims.UserID),
		)

		return nil, shortcut.ErrInvalidToken
	}

	s.logger.Info("checking redis session",
		zap.String("access_jti", accessTokenClaims.ID),
	)

	err = s.redis.CheckSessionStatus(ctx, accessTokenClaims.ID)
	if err != nil {
		s.logger.Info("session check failed",
			zap.String("access_jti", accessTokenClaims.ID),
			zap.Error(err),
		)

		if errors.Is(err, shortcut.ErrSessionNotFound) {
			s.logger.Info("session expired, starting refresh flow",
				zap.String("access_jti", accessTokenClaims.ID),
				zap.String("refresh_jti", refTokenClaims.ID),
				zap.String("user_id", refTokenClaims.UserID),
			)

			hashRefToken, err := s.authRepo.GetActiveRefreshToken(
				ctx,
				refTokenClaims.ID,
				refTokenClaims.UserID,
				userAgent,
			)
			if err != nil {
				s.logger.Error("failed to get active refresh token",
					zap.String("refresh_jti", refTokenClaims.ID),
					zap.String("user_id", refTokenClaims.UserID),
					zap.Error(err),
				)

				return nil, err
			}

			s.logger.Info("active refresh token found",
				zap.String("refresh_jti", refTokenClaims.ID),
			)

			err = utils.Compare(hashRefToken, tokens.RefreshToken)
			if err != nil {
				s.logger.Warn("refresh token hash mismatch",
					zap.String("refresh_jti", refTokenClaims.ID),
					zap.Error(err),
				)

				return nil, shortcut.ErrInvalidToken
			}

			s.logger.Info("refresh token validated",
				zap.String("refresh_jti", refTokenClaims.ID),
			)

			accessJTI := uuid.New().String()

			newAccessToken, err := utils.GenerateAccessToken(
				accessTokenClaims.UserID,
				s.secret,
				accessJTI,
			)
			if err != nil {
				s.logger.Error("failed to generate access token",
					zap.String("user_id", accessTokenClaims.UserID),
					zap.Error(err),
				)

				return nil, err
			}

			refreshJTI := uuid.New()

			newRefToken, err := utils.GenerateRefreshToken(
				refTokenClaims.UserID,
				s.secret,
				refreshJTI.String(),
			)
			if err != nil {
				s.logger.Error("failed to generate refresh token",
					zap.String("user_id", refTokenClaims.UserID),
					zap.Error(err),
				)

				return nil, err
			}

			s.logger.Info("new token pair generated",
				zap.String("new_access_jti", accessJTI),
				zap.String("new_refresh_jti", refreshJTI.String()),
			)

			hashNewRefToken, err := utils.Hash(newRefToken)
			if err != nil {
				s.logger.Error("failed to hash refresh token",
					zap.Error(err),
				)

				return nil, err
			}

			userID, _ := uuid.Parse(refTokenClaims.UserID)

			storageRefreshToken := &dto.RefreshToken{
				ID:        refreshJTI,
				UserID:    userID,
				TokenHash: hashNewRefToken,
				ExpiresAt: time.Now().Add(72 * time.Hour),
				RevokedAt: nil,
				CreatedAt: time.Now(),
				UserAgent: userAgent,
				IP:        ip,
			}

			oldRefTokenID, _ := uuid.Parse(refTokenClaims.ID)

			err = s.authRepo.UpdateRefreshTokenByID(
				ctx,
				oldRefTokenID,
				storageRefreshToken,
			)
			if err != nil {
				s.logger.Error("failed to rotate refresh token",
					zap.String("old_refresh_jti", oldRefTokenID.String()),
					zap.String("new_refresh_jti", refreshJTI.String()),
					zap.Error(err),
				)

				return nil, err
			}

			s.logger.Info("refresh token rotated",
				zap.String("old_refresh_jti", oldRefTokenID.String()),
				zap.String("new_refresh_jti", refreshJTI.String()),
			)

			session := s.redis.NewSession(accessTokenClaims.UserID)

			err = s.redis.SaveSession(
				ctx,
				accessJTI,
				session,
				time.Minute*15,
			)
			if err != nil {
				s.logger.Error("failed to save new session",
					zap.String("access_jti", accessJTI),
					zap.Error(err),
				)

				return nil, err
			}

			s.logger.Info("new session created",
				zap.String("access_jti", accessJTI),
				zap.String("user_id", accessTokenClaims.UserID),
			)

			return &dto.Tokens{
				RefreshToken: newRefToken,
				AccessToken:  newAccessToken,
			}, nil
		}

		s.logger.Warn("session validation failed",
			zap.String("access_jti", accessTokenClaims.ID),
			zap.Error(err),
		)

		return nil, shortcut.ErrInvalidToken
	}

	s.logger.Info("check passed",
		zap.String("access_jti", accessTokenClaims.ID),
		zap.String("user_id", accessTokenClaims.UserID),
	)

	return tokens, nil
}
