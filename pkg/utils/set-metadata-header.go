package utils

import (
	"auth-micro-service/internal/dto"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SetMetadataHeaderTokens(ctx context.Context, tokens *dto.Tokens) error {
	header := metadata.Pairs(
		"authorization", "Bearer "+tokens.AccessToken,
		"refresh-token", "Bearer "+tokens.RefreshToken,
	)

	if err := grpc.SetHeader(ctx, header); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
