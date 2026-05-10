package auth

import (
	"context"
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
		return "", err
	}

	return refreshTokenHash, nil
}
