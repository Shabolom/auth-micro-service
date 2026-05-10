package delet

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/render"
	"context"
)

type DeleteUsersService interface {
	DeleteUser(ctx context.Context, userID string) error
}

type Handler struct {
	deleteUsersService DeleteUsersService
}

func New(deleteUsersService DeleteUsersService) *Handler {
	return &Handler{
		deleteUsersService: deleteUsersService,
	}
}

func (h *Handler) DeleteUsers(ctx context.Context, req *authv1.DeleteUsersRequest) (*authv1.DeleteUsersReply, error) {
	userID := req.GetId()
	err := h.deleteUsersService.DeleteUser(ctx, userID)
	if err != nil {
		return nil, render.Error(err)
	}
	return &authv1.DeleteUsersReply{
		Message: "User deleted",
	}, nil
}
