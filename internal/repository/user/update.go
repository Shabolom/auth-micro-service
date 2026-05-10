package user

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/internal/dto"
	"context"
	"fmt"
)

func (s *Storage) UpdateUser(ctx context.Context, accountID string, req *authv1.UpdateUser) (*dto.AccountAndUser, error) {
	const query = `
		WITH updated_account AS (
			UPDATE accounts
			SET
				email = $2,
				updated_at = NOW()
			WHERE id = $1
			  AND deleted_at IS NULL
			RETURNING id, email
		),
		updated_user AS (
			UPDATE users
			SET
				name = $3,
				age = $4,
				updated_at = NOW()
			WHERE account_id = $1
			  AND deleted_at IS NULL
			RETURNING
				account_id,
				name,
				age,
				created_at,
				updated_at,
				deleted_at
		)
		SELECT
			a.id,
			a.email,
			u.name,
			u.age,
			u.created_at,
			u.updated_at,
			u.deleted_at
		FROM updated_account a
		JOIN updated_user u ON u.account_id = a.id
	`

	var res dto.AccountAndUser

	err := s.conn.QueryRow(
		ctx,
		query,
		accountID,
		req.GetMail(),
		req.GetName(),
		req.GetAge(),
	).Scan(
		&res.ID,
		&res.Email,
		&res.Name,
		&res.Age,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("UpdateUser query row: %w", err)
	}

	return &res, nil
}
