package list

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UsersListService interface {
	GetUsers(ctx context.Context) ([]*dto.AccountAndUser, error)
}

type Handler struct {
	usersListService UsersListService
}

func New(getUsersListService UsersListService) *Handler {
	return &Handler{
		usersListService: getUsersListService,
	}
}

func (h *Handler) GetUsersList(ctx context.Context, req *emptypb.Empty) (*authv1.GetUsersReply, error) {
	accounts, err := h.usersListService.GetUsers(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	users := make([]*authv1.User, 0, len(accounts))

	for _, account := range accounts {
		user := &authv1.User{
			Id:        account.ID,
			Mail:      account.Email,
			Name:      account.Name,
			Age:       uint32(account.Age),
			CreatedAt: timestamppb.New(account.CreatedAt),
			AddedAt:   timestamppb.New(account.UpdatedAt),
		}

		users = append(users, user)
	}

	return &authv1.GetUsersReply{
		ErrInfoReason: authv1.GetUsersReply_STATUS_OK,
		Users:         users,
	}, nil
}
