package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

	now := time.Now()

	accessTokenJTI := uuid.New()
	accessToken, err := utils.GenerateAccessToken(account.ID, s.secret, accessTokenJTI.String())
	if err != nil {
		s.logger.Info("Error generating access token")
		return &dto.Tokens{}, err
	}

	refreshTokenJTI := uuid.New()
	refreshToken, err := utils.GenerateRefreshToken(account.ID, s.secret, refreshTokenJTI.String())
	if err != nil {
		s.logger.Info("Error generating refresh token")
		return &dto.Tokens{}, err
	}

	hashToken, err := utils.Hash(refreshToken)
	if err != nil {
		s.logger.Info("Error hashing refresh token")
		return &dto.Tokens{}, err
	}

	userID, err := uuid.Parse(account.ID)
	if err != nil {
		s.logger.Info("Error parsing account id")
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
		s.logger.Info("Error creating refresh token")
		return &dto.Tokens{}, err
	}

	session := s.inmemorystorage.NewSession(account.ID)
	s.inmemorystorage.Save(accessTokenJTI.String(), session)

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
