package shortcut

import "errors"

var (
	ErrAccountNotFound              = errors.New("account not found")
	ErrEmailAlreadyExists           = errors.New("email already exists")
	ErrWrongPasswordOrEmail         = errors.New("wrong password or email")
	ErrInvalidTokenPair             = errors.New("access token and refresh token belong to different users")
	ErrRefreshSessionNotFound       = errors.New("refresh session not found")
	ErrInvalidToken                 = errors.New("invalid token")
	ErrRevokedSession               = errors.New("session revoked")
	ErrUserDescriptionAlreadyExists = errors.New("user description already exists")
)
