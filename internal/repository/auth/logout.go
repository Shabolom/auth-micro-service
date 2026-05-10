package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) Logout(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE refresh_sessions
		SET revoked_at = NOW()
		WHERE id = $1
		  AND revoked_at IS NULL
	`

	tag, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Logout: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("Logout: session not found")
	}

	return nil
}
