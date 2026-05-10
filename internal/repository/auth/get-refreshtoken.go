package auth

import (
	"auth-micro-service/pkg/shortcut"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetActiveRefreshToken(ctx context.Context, jti string, userID string, userAgent string) (string, error) {
	const query = `
		SELECT
			refresh_token_hash
		FROM refresh_sessions
		WHERE id = $1
		  AND user_id = $2
		  AND user_agent = $3
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
	`

	var refreshTokenHash string

	err := s.conn.QueryRow(
		ctx,
		query,
		jti,
		userID,
		userAgent,
	).Scan(
		&refreshTokenHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", shortcut.ErrNoRows
		}
		return "", err
	}

	return refreshTokenHash, nil
}
