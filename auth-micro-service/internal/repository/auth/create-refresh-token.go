package auth

import (
	"auth-micro-service/internal/dto"
	"context"
)

func (s *Storage) CreateRefreshToken(ctx context.Context, session *dto.RefreshToken) error {
	const query = `
		INSERT INTO refresh_sessions (
			id,
			user_id,
			refresh_token_hash,
			expires_at,
			user_agent,
			ip
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
	`

	_, err := s.conn.Exec(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.ExpiresAt,
		session.UserAgent,
		session.IP,
	)
	if err != nil {
		return err
	}

	return nil
}
