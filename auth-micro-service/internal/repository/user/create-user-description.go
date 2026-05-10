package user

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/shortcut"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateUserDescription(ctx context.Context, accountID string, name *string, age *int) (*dto.UserDescription, error) {
	const query = `
		INSERT INTO users (account_id,name,age)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			account_id,
			name,
			age,
			created_at,
			updated_at,
			deleted_at
	`

	var user dto.UserDescription

	err := s.conn.QueryRow(
		ctx,
		query,
		accountID,
		name,
		age,
	).Scan(
		&user.AccountID,
		&user.Name,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, shortcut.ErrUserDescriptionAlreadyExists
			}
		}

		return nil, err
	}

	return &user, nil
}
