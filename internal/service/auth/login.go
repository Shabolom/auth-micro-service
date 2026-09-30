package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/rabbitMQ"
	"auth-micro-service/pkg/shortcut"
	"auth-micro-service/pkg/utils"
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) Login(ctx context.Context, login *dto.LoginRequest) (*dto.Tokens, error) {
	s.logger.Info(
		"login started",
		zap.String("email", login.Email),
		zap.String("ip", login.IP),
		zap.String("user_agent", login.UserAgent),
	)

	err := s.logiValidate(login)
	if err != nil {
		s.logger.Warn(
			"login validation failed",
			zap.String("email", login.Email),
			zap.Error(err),
		)
		return &dto.Tokens{}, err
	}

	account, err := s.authRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		s.logger.Warn(
			"login failed: wrong email or account repo error",
			zap.String("email", login.Email),
			zap.Error(err),
		)
		return &dto.Tokens{}, shortcut.ErrWrongPasswordOrEmail
	}

	if err = utils.Compare(account.PasswordHash, login.Password); err != nil {
		s.logger.Warn(
			"login failed: wrong password",
			zap.String("email", login.Email),
			zap.String("user_id", account.ID),
			zap.Error(err),
		)
		return &dto.Tokens{}, shortcut.ErrWrongPasswordOrEmail
	}

	accessTokenJTI := uuid.New()
	s.logger.Debug(
		"login access token jti generated",
		zap.String("email", login.Email),
		zap.String("user_id", account.ID),
		zap.String("access_token_jti", accessTokenJTI.String()),
	)

	accessToken, err := utils.GenerateAccessToken(account.ID, s.secret, accessTokenJTI.String())
	if err != nil {
		s.logger.Error(
			"login failed: generate access token",
			zap.String("email", login.Email),
			zap.String("user_id", account.ID),
			zap.Error(err),
		)
		return &dto.Tokens{}, err
	}

	refreshTokenJTI := uuid.New()
	s.logger.Debug(
		"login refresh token jti generated",
		zap.String("email", login.Email),
		zap.String("user_id", account.ID),
		zap.String("refresh_token_jti", refreshTokenJTI.String()),
	)

	refreshToken, err := utils.GenerateRefreshToken(account.ID, s.secret, refreshTokenJTI.String())
	if err != nil {
		s.logger.Error(
			"login failed: generate refresh token",
			zap.String("email", login.Email),
			zap.String("user_id", account.ID),
			zap.Error(err),
		)
		return &dto.Tokens{}, err
	}

	hashToken, err := utils.Hash(refreshToken)
	if err != nil {
		s.logger.Error(
			"login failed: hash refresh token",
			zap.String("email", login.Email),
			zap.String("user_id", account.ID),
			zap.Error(err),
		)
		return &dto.Tokens{}, err
	}

	userID, err := uuid.Parse(account.ID)
	if err != nil {
		s.logger.Error(
			"login failed: parse account id",
			zap.String("email", login.Email),
			zap.String("user_id", account.ID),
			zap.Error(err),
		)
		return &dto.Tokens{}, err
	}

	repoRefreshToken := &dto.RefreshToken{
		ID:        refreshTokenJTI,
		UserID:    userID,
		TokenHash: hashToken,
		ExpiresAt: time.Now().Add(72 * time.Hour),
		RevokedAt: nil,
		CreatedAt: time.Now(),
		IP:        login.IP,
		UserAgent: login.UserAgent,
	}

	s.logger.Debug(
		"login refresh token prepared",
		zap.String("email", login.Email),
		zap.String("user_id", account.ID),
		zap.String("refresh_token_jti", refreshTokenJTI.String()),
		zap.Time("expires_at", repoRefreshToken.ExpiresAt),
	)

	if err = s.authRepo.CreateRefreshToken(ctx, repoRefreshToken); err != nil {
		s.logger.Error(
			"login failed: create refresh token",
			zap.String("email", login.Email),
			zap.String("user_id", account.ID),
			zap.String("refresh_token_jti", refreshTokenJTI.String()),
			zap.Error(err),
		)
		return &dto.Tokens{}, err
	}

	session := s.redis.NewSession(account.ID)

	if err = s.redis.SaveSession(ctx, accessTokenJTI.String(), session, time.Minute*15); err != nil {
		s.logger.Error(
			"login failed: save redis session",
			zap.String("email", login.Email),
			zap.String("user_id", account.ID),
			zap.String("access_token_jti", accessTokenJTI.String()),
			zap.Duration("ttl", time.Minute*15),
			zap.Error(err),
		)
		return &dto.Tokens{}, err
	}

	s.logger.Info(
		"login session saved",
		zap.String("email", login.Email),
		zap.String("user_id", account.ID),
		zap.String("access_token_jti", accessTokenJTI.String()),
		zap.String("refresh_token_jti", refreshTokenJTI.String()),
	)

	go func(email string) {
		publishCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		s.logger.Debug(
			"login event publishing started",
			zap.String("email", email),
			zap.String("routing_key", "login"),
		)

		err := s.rabbitMQ.Publish(publishCtx, "login", rabbitMQ.TEXTTYPE, []byte(email))
		if err != nil {
			s.logger.Warn(
				"login event publishing failed",
				zap.String("email", email),
				zap.String("routing_key", "login"),
				zap.Error(err),
			)
			return
		}

		s.logger.Debug(
			"login event published",
			zap.String("email", email),
			zap.String("routing_key", "login"),
		)
	}(login.Email)

	s.logger.Info(
		"login completed",
		zap.String("email", login.Email),
		zap.String("user_id", account.ID),
	)

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) logiValidate(req *dto.LoginRequest) error {
	if req.Email == "" {
		s.logger.Warn("login email is empty")
		return shortcut.ErrEmptyCredentials
	}

	if req.Password == "" {
		s.logger.Warn(
			"login password is empty",
			zap.String("email", req.Email),
		)
		return shortcut.ErrEmptyCredentials
	}

	if len(req.Password) < 5 {
		s.logger.Warn(
			"login password is too short",
			zap.String("email", req.Email),
			zap.Int("password_len", len(req.Password)),
		)
		return shortcut.ErrEmptyCredentials
	}

	err := utils.ValidateEmail(req.Email)
	if err != nil {
		s.logger.Warn(
			"login email validation failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return shortcut.ErrValidateEmail
	}

	return nil
}
