package user

import (
	"auth-micro-service/internal/dto"
	"context"
	"fmt"
)

func (s *Storage) GetAccountByID(ctx context.Context, accountID string) (*dto.AccountAndUser, error) {
	const query = `
		SELECT
			a.id,
			a.email,
			a.created_at,
			a.updated_at,
			a.deleted_at,

			u.name,
			u.age
		FROM accounts a
		LEFT JOIN users u
			ON u.account_id = a.id
			AND u.deleted_at IS NULL
		WHERE a.id = $1
		  AND a.deleted_at IS NULL
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
		return &dto.AccountAndUser{}, fmt.Errorf("GetAccountByID: %w", err)
	}

	return &account, nil
}
