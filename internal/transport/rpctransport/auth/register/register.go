package register_handler

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"auth-micro-service/internal/render"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type RegisterService interface {
	Register(ctx context.Context, request *dto.RegisterRequest) (dto.Tokens, error)
}

type Handler struct {
	registerService RegisterService
}

func New(registerService RegisterService) *Handler {
	return &Handler{
		registerService: registerService,
	}
}

func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterReply, error) {
	ip, userAgent := getRequestInfo(ctx)

	registerRequest := &dto.RegisterRequest{
		Email:     req.GetMail(),
		Password:  req.GetPassword(),
		IP:        ip,
		UserAgent: userAgent,
		Name:      req.GetName(),
		Age:       int(req.GetAge()),
	}

	tokens, err := h.registerService.Register(ctx, registerRequest)
	if err != nil {
		return nil, render.Error(err)
	}

	if err = setTokenHeaders(ctx, tokens); err != nil {
		log.Println(err)
		return nil, render.Error(err)
	}

	return &authv1.RegisterReply{
		ErrInfoReason: authv1.RegisterReply_STATUS_OK,
		Message:       "registered success",
	}, nil
}

func getRequestInfo(ctx context.Context) (ip string, userAgent string) {
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

	return ip, userAgent
}

func setTokenHeaders(ctx context.Context, tokens dto.Tokens) error {
	header := metadata.Pairs(
		"authorization", "Bearer "+tokens.AccessToken,
		"refresh-token", "Bearer "+tokens.RefreshToken,
	)

	return grpc.SetHeader(ctx, header)
}
