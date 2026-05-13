package redis

import (
	"context"

	"go.uber.org/zap"
)

func (r *Redis) RevokeSession(ctx context.Context, key string) error {
	err := r.client.HSet(ctx, key, map[string]interface{}{
		REVOKE: true,
	}).Err()

	if err != nil {
		r.logger.Error("redis HSet", zap.Error(err))
		return err
	}

	return nil
}
