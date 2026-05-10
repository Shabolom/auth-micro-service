package utils

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

func AccessTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found")
	}
	
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", errors.New("authorization header not found")
	}
	token := authHeaders[0]

	return strings.TrimPrefix(token, "Bearer "), nil
}

func RefreshTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found")
	}

	authHeaders := md.Get("refresh-token")
	if len(authHeaders) == 0 {
		return "", errors.New("authorization header not found")
	}
	token := authHeaders[0]

	return strings.TrimPrefix(token, "Bearer "), nil
}
