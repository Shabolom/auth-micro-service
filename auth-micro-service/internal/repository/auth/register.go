package auth

import (
	"auth-micro-service/internal/dto"
	"auth-micro-service/pkg/shortcut"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) Register(ctx context.Context, reg *dto.Register) error {
	const query = `
		INSERT INTO accounts (
			id,
			email,
			password_hash
		)
		VALUES (
			$1,
			$2,
			$3
		)
	`

	_, err := s.conn.Exec(ctx, query,
		reg.ID,
		reg.Email,
		reg.PasswordHash,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "ux_accounts_email" {
				return shortcut.ErrEmailAlreadyExists
			}
		}

		return fmt.Errorf("Register account: %w", err)
	}

	return nil
}
