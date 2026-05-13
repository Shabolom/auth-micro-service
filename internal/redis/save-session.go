package redis

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func (r *Redis) SaveSession(ctx context.Context, key string, value *Session, expiration time.Duration) error {
	pipe := r.client.TxPipeline()

	pipe.HSet(ctx, key, map[string]interface{}{
		USERID: value.UserID,
		EXPIRE: value.ExpiresAt,
		REVOKE: value.Revoked,
	})

	pipe.Expire(ctx, key, expiration)

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("redis save-session", zap.Error(err))
		return err
	}

	return nil
}
