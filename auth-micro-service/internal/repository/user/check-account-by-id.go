package user

import (
	"auth-micro-service/pkg/shortcut"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) CheckAccountByID(ctx context.Context, userID string) error {
	const query = `
		SELECT 1
		FROM accounts
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	var exists int

	err := s.conn.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shortcut.ErrAccountNotFound
		}
		return err
	}

	return nil
}
