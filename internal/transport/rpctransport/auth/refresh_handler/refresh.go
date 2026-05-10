package refresh_handler

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"auth-micro-service/pkg/utils"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
	ip := ""
	userAgent := ""

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		agents := md.Get("user-agent")
		if len(agents) > 0 {
			userAgent = agents[0]
		}
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		ip = p.Addr.String()
	}

	oldRefreshToken, err := utils.RefreshTokenFromMetadata(ctx)
	if err != nil {
		return nil, render.Error(err)
	}

	tokens, err := s.refreshService.Refresh(ctx, oldRefreshToken, userAgent, ip)
	if err != nil {
		return nil, render.Error(err)
	}

	header := metadata.Pairs(
		"authorization", "Bearer "+tokens.AccessToken,
		"refresh-token", "Bearer "+tokens.RefreshToken,
	)

	if err = grpc.SetHeader(ctx, header); err != nil {
		log.Println(err)
		return nil, render.Error(err)
	}

	return &authv1.RefreshReply{
		Message: "success",
	}, nil
}
