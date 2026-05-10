package auth

import (
	"auth-micro-service/internal/dto"
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrRefreshSessionNotFound = errors.New("refresh session not found")

func (s *Storage) UpdateRefreshTokenByID(ctx context.Context, oldJTI uuid.UUID, session *dto.RefreshToken) error {
	const query = `
		UPDATE refresh_sessions
		SET
			id = $1,
			refresh_token_hash = $2,
			expires_at = $3,
			user_agent = $4,
			ip = $5,
			revoked_at = NULL
		WHERE id = $6
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
	`

	tag, err := s.conn.Exec(ctx, query,
		session.ID,
		session.TokenHash,
		session.ExpiresAt,
		session.UserAgent,
		session.IP,
		oldJTI,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrRefreshSessionNotFound
	}

	return nil
}
