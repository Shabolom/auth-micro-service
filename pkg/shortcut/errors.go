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
	ErrEmptyFields                  = errors.New("empty field is not allowed")
	ErrAgeLimit                     = errors.New("age is too young")
	ErrEmptyCredentials             = errors.New("credentials is empty")
	ErrNoRows                       = errors.New("the record was not found")
	ErrValidateEmail                = errors.New("the email is invalid")
)
