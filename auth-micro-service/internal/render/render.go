package render

import (
	"auth-micro-service/pkg/shortcut"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Error(err error) error {
	switch {

	case errors.Is(err, shortcut.ErrAccountNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, shortcut.ErrEmailAlreadyExists),
		errors.Is(err, shortcut.ErrUserDescriptionAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, shortcut.ErrWrongPasswordOrEmail):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, shortcut.ErrInvalidToken),
		errors.Is(err, shortcut.ErrInvalidTokenPair),
		errors.Is(err, shortcut.ErrRefreshSessionNotFound),
		errors.Is(err, shortcut.ErrRevokedSession):
		return status.Error(codes.Unauthenticated, err.Error())

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
