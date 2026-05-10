package update

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/utils"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UpdateUsersService interface {
	UpdateUsers(ctx context.Context, strAccessToken string, req *authv1.UpdateUser) (*dto.AccountAndUser, error)
}

type Handler struct {
	updateUsersService UpdateUsersService
}

func New(updateUsersService UpdateUsersService) *Handler {
	return &Handler{
		updateUsersService: updateUsersService,
	}
}

func (h *Handler) UpdateUsers(ctx context.Context, req *authv1.UpdateUsersRequest) (*authv1.UpdateUsersReply, error) {
	updatedUser := req.GetUpdatedUser()

	accessToken, err := utils.AccessTokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	userAndAccount, err := h.updateUsersService.UpdateUsers(ctx, accessToken, updatedUser)
	if err != nil {
		return nil, err
	}

	user := &authv1.User{
		Id:        userAndAccount.ID,
		Mail:      userAndAccount.Email,
		Name:      userAndAccount.Name,
		Age:       uint32(userAndAccount.Age),
		CreatedAt: timestamppb.New(userAndAccount.CreatedAt),
		AddedAt:   timestamppb.New(userAndAccount.UpdatedAt),
	}

	return &authv1.UpdateUsersReply{
		ErrInfoReason: authv1.UpdateUsersReply_STATUS_OK,
		User:          user,
		Message:       "success",
	}, nil
}
