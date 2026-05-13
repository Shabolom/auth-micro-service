package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/shortcut"
	"auth-micro-service/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

const refreshTokenTTL = 72 * time.Hour

var now = time.Now()

func (s *Service) Register(ctx context.Context, request *dto.RegisterRequest) (dto.Tokens, error) {
	err := requestValidate(request)
	if err != nil {
		s.logger.Info("validation error", zap.Error(err))
		return dto.Tokens{}, err
	}

	hashPassword, err := utils.Hash(request.Password)
	if err != nil {
		s.logger.Info("Error hashing password")
		return dto.Tokens{}, err
	}

	userID := uuid.New()

	register := &dto.Register{
		ID:           userID.String(),
		Email:        request.Email,
		PasswordHash: hashPassword,
		Name:         request.Name,
		Age:          request.Age,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err = s.authRepo.Register(ctx, register); err != nil {
		if errors.Is(err, shortcut.ErrEmailAlreadyExists) {
			s.logger.Info("email already exists", zap.Error(err))
			return dto.Tokens{}, err
		}

		s.logger.Info("Error registering user", zap.Error(err))
		return dto.Tokens{}, err
	}

	accessTokenJTI := uuid.New()
	accessToken, err := utils.GenerateAccessToken(userID.String(), s.secret, accessTokenJTI.String())
	if err != nil {
		s.logger.Info("Error generating access token", zap.Error(err))
		return dto.Tokens{}, err
	}

	refreshTokenJTI := uuid.New()
	refreshToken, err := utils.GenerateRefreshToken(userID.String(), s.secret, refreshTokenJTI.String())
	if err != nil {
		s.logger.Warn("Error generating refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	hashToken, err := utils.Hash(refreshToken)
	if err != nil {
		s.logger.Warn("Error hashing refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	repoRefreshToken := &dto.RefreshToken{
		ID:        refreshTokenJTI,
		UserID:    userID,
		TokenHash: hashToken,
		ExpiresAt: now.Add(refreshTokenTTL),
		RevokedAt: nil,
		CreatedAt: now,
		UserAgent: request.UserAgent,
		IP:        request.IP,
	}

	if err = s.authRepo.CreateRefreshToken(ctx, repoRefreshToken); err != nil {
		s.logger.Warn("Error creating refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	session := s.redis.NewSession(userID.String())
	err = s.redis.SaveSession(ctx, accessTokenJTI.String(), session, time.Minute*15)
	if err != nil {
		s.logger.Info("Error redis-storage save session", zap.Error(err))
		return dto.Tokens{}, err
	}

	go func(email string) {
		publishCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := s.rabbitMQ.Publish(publishCtx, "register", shortcut.TEXTTYPE, []byte(email))
		if err != nil {
			log.Error("Error publishing email", zap.Error(err))
		}
	}(register.Email)

	return dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func requestValidate(req *dto.RegisterRequest) error {

	if req.Email == "" || req.Password == "" {
		return shortcut.ErrEmptyCredentials
	}
	if req.Name == "" || req.Age <= 0 {
		return shortcut.ErrEmptyFields
	}

	return nil
}
