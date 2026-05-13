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

func (s *Service) Login(ctx context.Context, login *dto.LoginRequest) (*dto.Tokens, error) {
	if login.Email == "" || login.Password == "" {
		s.logger.Warn("login email or password is empty")
		return &dto.Tokens{}, errors.New("email or password is empty")
	}

	account, err := s.authRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		s.logger.Info(err.Error())
		return &dto.Tokens{}, errors.New("wrong password or email")
	}

	if err = utils.Compare(account.PasswordHash, login.Password); err != nil {
		s.logger.Info(err.Error())
		return &dto.Tokens{}, errors.New("wrong password or email")
	}

	accessTokenJTI := uuid.New()
	accessToken, err := utils.GenerateAccessToken(account.ID, s.secret, accessTokenJTI.String())
	if err != nil {
		s.logger.Info("Error generating access token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	refreshTokenJTI := uuid.New()
	refreshToken, err := utils.GenerateRefreshToken(account.ID, s.secret, refreshTokenJTI.String())
	if err != nil {
		s.logger.Info("Error generating refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	hashToken, err := utils.Hash(refreshToken)
	if err != nil {
		s.logger.Info("Error hashing refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	userID, err := uuid.Parse(account.ID)
	if err != nil {
		s.logger.Info("Error parsing account id", zap.Error(err))
		return &dto.Tokens{}, err
	}

	repoRefreshToken := &dto.RefreshToken{
		ID:        refreshTokenJTI,
		UserID:    userID,
		TokenHash: hashToken,
		ExpiresAt: now.Add(72 * time.Hour),
		RevokedAt: nil,
		CreatedAt: now,
		IP:        login.IP,
		UserAgent: login.UserAgent,
	}

	if err = s.authRepo.CreateRefreshToken(ctx, repoRefreshToken); err != nil {
		s.logger.Info("Error creating refresh token", zap.Error(err))
		return &dto.Tokens{}, err
	}

	session := s.redis.NewSession(account.ID)
	err = s.redis.SaveSession(ctx, accessTokenJTI.String(), session, time.Minute*15)
	if err != nil {
		s.logger.Info("Error creating redis-storage session", zap.Error(err))
		return &dto.Tokens{}, err
	}

	go func(email string) {
		publishCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := s.rabbitMQ.Publish(publishCtx, "login", shortcut.TEXTTYPE, []byte(email))
		if err != nil {
			s.logger.Info("Error publishing login", zap.Error(err))
		}
	}(login.Email)

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
