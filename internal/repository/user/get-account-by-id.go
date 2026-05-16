package user

import (
	"auth-micro-service/internal/dto"
	"context"
	"fmt"
)

func (s *Storage) GetUserByID(ctx context.Context, accountID string) (*dto.AccountAndUser, error) {
	const query = `
		SELECT
			id,
			email,
			created_at,
			updated_at,
			deleted_at,
			name,
			age
		FROM accounts
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	var account dto.AccountAndUser

	err := s.conn.QueryRow(ctx, query, accountID).Scan(
		&account.ID,
		&account.Email,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.DeletedAt,
		&account.Name,
		&account.Age,
	)
	if err != nil {
		return &dto.AccountAndUser{}, fmt.Errorf("GetUserByID: %w", err)
	}

	return &account, nil
}
