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
	GetAccounts(ctx context.Context) ([]dto.AccountAndUser, error)
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
	accounts, err := h.usersListService.GetAccounts(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	users := make([]*authv1.User, 0, len(accounts))

	for _, account := range accounts {
		user := &authv1.User{
			Id:   account.ID,
			Mail: account.Email,
		}

		if account.CreatedAt != nil {
			user.CreatedAt = timestamppb.New(*account.CreatedAt)
		}

		if account.UpdatedAt != nil {
			user.AddedAt = timestamppb.New(*account.UpdatedAt)
		}

		if account.Name != nil {
			user.Name = *account.Name
		}

		if account.Age != nil {
			user.Age = uint32(*account.Age)
		}

		users = append(users, user)
	}

	return &authv1.GetUsersReply{
		Users: users,
	}, nil
}
