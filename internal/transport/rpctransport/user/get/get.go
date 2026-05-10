package get

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetUsersService interface {
	GetUserByAccountID(ctx context.Context, accountID string) (*dto.AccountAndUser, error)
}

type Handler struct {
	getUsersService GetUsersService
}

func New(getUsersService GetUsersService) *Handler {
	return &Handler{
		getUsersService: getUsersService,
	}
}

func (h *Handler) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserReply, error) {
	accountID := req.GetId()
	if accountID == "" {
		return nil, errors.New("accountID is empty")
	}

	userAccount, err := h.getUsersService.GetUserByAccountID(ctx, accountID)
	if err != nil {
		return nil, render.Error(err)
	}

	user := &authv1.User{
		Id:   userAccount.ID,
		Mail: userAccount.Email,
	}

	if userAccount.CreatedAt != nil {
		user.CreatedAt = timestamppb.New(*userAccount.CreatedAt)
	}

	if userAccount.UpdatedAt != nil {
		user.AddedAt = timestamppb.New(*userAccount.UpdatedAt)
	}

	if userAccount.Name != nil {
		user.Name = *userAccount.Name
	}

	if userAccount.Age != nil {
		user.Age = uint32(*userAccount.Age)
	}

	return &authv1.GetUserReply{
		User: user,
	}, nil
}
