package refresh_handler

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"auth-micro-service/pkg/utils"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type RefreshService interface {
	Refresh(ctx context.Context, oldRefToken string, userAgent string, ip string) (*dto.Tokens, error)
}

type Handler struct {
	refreshService RefreshService
}

func New(refreshService RefreshService) *Handler {
	return &Handler{
		refreshService: refreshService,
	}
}
func (s *Handler) Refresh(ctx context.Context, req *emptypb.Empty) (*authv1.RefreshReply, error) {
	ip, userAgent := utils.IpUserAgentFromMetadata(ctx)

	oldRefreshToken, err := utils.RefreshTokenFromMetadata(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	tokens, err := s.refreshService.Refresh(ctx, oldRefreshToken, userAgent, ip)
	if err != nil {
		return nil, render.Error(err)
	}

	err = utils.SetMetadataHeaderTokens(ctx, tokens)
	if err != nil {
		return nil, render.Error(err)
	}

	return &authv1.RefreshReply{
		ErrInfoReason: authv1.RefreshReply_STATUS_OK,
		Message:       "success",
	}, nil
}
