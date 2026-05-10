package user

import (
	"auth-micro-service/internal/dto"
	"context"
	"fmt"
)

func (s *Storage) GetAccounts(ctx context.Context) ([]dto.AccountAndUser, error) {
	const query = `
		SELECT
			a.id,
			a.email,
			a.created_at,
			a.updated_at,
			a.deleted_at,

			u.name,
			u.age,
			u.created_at AS user_created_at,
			u.updated_at AS user_updated_at,
			u.deleted_at AS user_deleted_at
		FROM accounts a
		LEFT JOIN users u
			ON u.account_id = a.id
			AND u.deleted_at IS NULL
		WHERE a.deleted_at IS NULL
		ORDER BY a.created_at DESC
	`

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetAccounts query: %w", err)
	}
	defer rows.Close()

	accounts := make([]dto.AccountAndUser, 0)

	for rows.Next() {
		var account dto.AccountAndUser

		err = rows.Scan(
			&account.ID,
			&account.Email,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.DeletedAt,

			&account.Name,
			&account.Age,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAccounts scan: %w", err)
		}

		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAccounts rows: %w", err)
	}

	return accounts, nil
}
