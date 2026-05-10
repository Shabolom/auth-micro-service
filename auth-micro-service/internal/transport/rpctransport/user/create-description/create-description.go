package create_description

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"auth-micro-service/pkg/utils"
	"context"
	"fmt"
)

type UserDescriptionService interface {
	Create(ctx context.Context, request *dto.CreateUserDescriptionRequest, accessToken string) (*dto.UserDescription, error)
}

type Handler struct {
	userDescriptionService UserDescriptionService
}

func New(userDescriptionService UserDescriptionService) *Handler {
	return &Handler{
		userDescriptionService: userDescriptionService,
	}
}

func (h *Handler) CreateUserDescription(ctx context.Context, req *authv1.CreateUserDescriptionRequest) (*authv1.CreateUserDescriptionReply, error) {
	userDescription := &dto.CreateUserDescriptionRequest{
		Name: req.GetName(),
		Age:  int(req.GetAge()),
	}

	fmt.Println(userDescription, 123123123123)
	accessToken, err := utils.AccessTokenFromMetadata(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	_, err = h.userDescriptionService.Create(ctx, userDescription, accessToken)
	if err != nil {
		return nil, render.Error(err)
	}

	return &authv1.CreateUserDescriptionReply{
		Message: "description added successfully",
	}, nil
}
