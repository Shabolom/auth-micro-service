package user

import (
	"context"
	"fmt"
)

func (s *Storage) DeleteUser(ctx context.Context, accountID string) error {
	const query = `
		UPDATE accounts
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	tag, err := s.conn.Exec(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteUser: user not found")
	}

	return nil
}
