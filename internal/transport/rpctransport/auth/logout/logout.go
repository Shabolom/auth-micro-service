package logout_handeler

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"auth-micro-service/pkg/utils"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LogoutService interface {
	Logout(ctx context.Context, userAgent string, tokens *dto.Tokens) error
}

type Handler struct {
	logoutService LogoutService
}

func New(logoutService LogoutService) *Handler {
	return &Handler{
		logoutService: logoutService,
	}
}

func (h *Handler) Logout(ctx context.Context, req *emptypb.Empty) (*authv1.LogoutReply, error) {
	_, userAgent := utils.IpUserAgentFromMetadata(ctx)

	refToken, err := utils.RefreshTokenFromMetadata(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	accessToken, err := utils.AccessTokenFromMetadata(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	err = h.logoutService.Logout(ctx, userAgent, &dto.Tokens{
		RefreshToken: refToken,
		AccessToken:  accessToken,
	},
	)

	if err != nil {
		return nil, render.Error(err)
	}

	return &authv1.LogoutReply{
		ErrInfoReason: authv1.LogoutReply_STATUS_OK,
		Message:       "success",
	}, nil
}
