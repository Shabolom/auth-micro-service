package auth

import (
	"auth-micro-service/internal/dto"
	"context"
)

func (s *Storage) GetByEmail(ctx context.Context, email string) (*dto.Register, error) {
	const query = `
		SELECT
			id,
			email,
			password_hash,
			created_at,
			updated_at,
			deleted_at
		FROM accounts
		WHERE email = $1
		  AND deleted_at IS NULL
	`

	var account dto.Register

	err := s.conn.QueryRow(ctx, query, email).Scan(
		&account.ID,
		&account.Email,
		&account.PasswordHash,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
