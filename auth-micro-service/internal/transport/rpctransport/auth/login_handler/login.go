package login_handler

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
)

type LoginService interface {
	Login(ctx context.Context, login *dto.LoginRequest) (*dto.Tokens, error)
	UpdateRefreshToken(ctx context.Context, oldTokens *dto.Tokens, userAgent string, ip string) (dto.Tokens, error)
}

type Handler struct {
	loginService LoginService
}

func New(loginService LoginService) *Handler {
	return &Handler{
		loginService: loginService,
	}
}

func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginReply, error) {
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

	oldAccessToken, _ := utils.AccessTokenFromMetadata(ctx)
	oldRefreshToken, _ := utils.RefreshTokenFromMetadata(ctx)

	if oldRefreshToken != "" && oldAccessToken != "" {
		oldTokens := &dto.Tokens{
			AccessToken:  oldAccessToken,
			RefreshToken: oldRefreshToken,
		}

		tokens, err := h.loginService.UpdateRefreshToken(ctx, oldTokens, userAgent, ip)
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

		return &authv1.LoginReply{
			Message: "success",
		}, nil
	}

	login := dto.LoginRequest{
		Email:     req.GetMail(),
		Password:  req.GetPassword(),
		IP:        ip,
		UserAgent: userAgent,
	}

	tokens, err := h.loginService.Login(ctx, &login)
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

	return &authv1.LoginReply{
		Message: "success",
	}, nil
}
