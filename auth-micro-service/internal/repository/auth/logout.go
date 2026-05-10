package auth

import (
	"context"

	"github.com/google/uuid"
)

func (s *Storage) Logout(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE refresh_sessions
		SET revoked_at = NOW()
		WHERE id = $1
		  AND revoked_at IS NULL
	`

	_, err := s.conn.Exec(ctx, query, id)
	return err
}
