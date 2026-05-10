package get

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"auth-micro-service/pkg/utils"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetUsersService interface {
	GetUserByID(ctx context.Context, strAccessToken string) (*dto.AccountAndUser, error)
}

type Handler struct {
	getUsersService GetUsersService
}

func New(getUsersService GetUsersService) *Handler {
	return &Handler{
		getUsersService: getUsersService,
	}
}

func (h *Handler) GetUser(ctx context.Context, req *emptypb.Empty) (*authv1.GetUserReply, error) {
	accessToken, err := utils.AccessTokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	userAccount, err := h.getUsersService.GetUserByID(ctx, accessToken)
	if err != nil {
		return nil, render.Error(err)
	}

	user := &authv1.User{
		Id:        userAccount.ID,
		Mail:      userAccount.Email,
		Name:      userAccount.Name,
		Age:       uint32(userAccount.Age),
		CreatedAt: timestamppb.New(userAccount.CreatedAt),
		AddedAt:   timestamppb.New(userAccount.UpdatedAt),
	}

	return &authv1.GetUserReply{
		ErrInfoReason: authv1.GetUserReply_STATUS_OK,
		User:          user,
	}, nil
}
