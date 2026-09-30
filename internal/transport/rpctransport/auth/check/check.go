package check

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"auth-micro-service/pkg/shortcut"
	"auth-micro-service/pkg/utils"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CheckService interface {
	Check(ctx context.Context, tokens *dto.Tokens, userAgent string, ip string) (*dto.Tokens, error)
}

type Handler struct {
	checkService CheckService
}

func New(checkService CheckService) *Handler {
	return &Handler{
		checkService: checkService,
	}
}

func (h *Handler) Check(ctx context.Context, req *emptypb.Empty) (*authv1.CheckReply, error) {
	accessToken, err := utils.AccessTokenFromMetadata(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	refreshToken, err := utils.RefreshTokenFromMetadata(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	if accessToken == "" || refreshToken == "" {
		return nil, render.Error(shortcut.ErrInvalidToken)
	}

	tokens := &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ip, userAgent := utils.IpUserAgentFromMetadata(ctx)

	newTokens, err := h.checkService.Check(ctx, tokens, userAgent, ip)
	if err != nil {
		return nil, render.Error(err)
	}

	if newTokens != nil {
		err = utils.SetMetadataHeaderTokens(ctx, newTokens)
		if err != nil {
			return nil, render.Error(err)
		}
	}

	return &authv1.CheckReply{
		ErrInfoReason: authv1.CheckReply_STATUS_OK,
		Message:       "success",
	}, nil
}
