package user

import (
	"auth-micro-service/internal/dto"
	"context"
	"fmt"
)

func (s *Storage) GetUsers(ctx context.Context) ([]*dto.AccountAndUser, error) {
	const query = `
		SELECT
			id,
			email,
			name,
			age,
			created_at,
			updated_at,
			deleted_at
		FROM accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetUsers query: %w", err)
	}
	defer rows.Close()

	accounts := make([]*dto.AccountAndUser, 0)

	for rows.Next() {
		var account dto.AccountAndUser

		err = rows.Scan(
			&account.ID,
			&account.Email,
			&account.Name,
			&account.Age,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetUsers scan: %w", err)
		}

		accounts = append(accounts, &account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUsers rows: %w", err)
	}

	return accounts, nil
}
