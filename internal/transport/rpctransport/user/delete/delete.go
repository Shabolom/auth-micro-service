package delet

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"auth-micro-service/pkg/utils"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type DeleteUsersService interface {
	DeleteUser(ctx context.Context, tokens *dto.Tokens) error
}

type Handler struct {
	deleteUsersService DeleteUsersService
}

func New(deleteUsersService DeleteUsersService) *Handler {
	return &Handler{
		deleteUsersService: deleteUsersService,
	}
}

func (h *Handler) DeleteUsers(ctx context.Context, req *emptypb.Empty) (*authv1.DeleteUsersReply, error) {
	strAccessToken, err := utils.AccessTokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	strRefreshToken, err := utils.RefreshTokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	tokens := &dto.Tokens{
		AccessToken:  strAccessToken,
		RefreshToken: strRefreshToken,
	}

	err = h.deleteUsersService.DeleteUser(ctx, tokens)
	if err != nil {
		return nil, render.Error(err)
	}
	return &authv1.DeleteUsersReply{
		ErrInfoReason: authv1.DeleteUsersReply_STATUS_OK,
		Message:       "User deleted",
	}, nil
}
