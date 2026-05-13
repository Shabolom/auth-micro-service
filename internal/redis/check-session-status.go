package redis

import (
	"context"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (r *Redis) CheckSessionStatus(ctx context.Context, jti string) error {

	result, err := r.client.HGet(ctx, jti, REVOKE).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Info("session status not exist", zap.String("jti", jti))
			return nil
		}

		return err
	}

	status, err := strconv.ParseBool(result)
	if err != nil {
		r.logger.Info("session status not bool", zap.String("jti", jti))
		return err
	}

	if status {
		r.logger.Info("session status ok", zap.String("jti", jti))
		return errors.New("session revoked")
	}

	return nil
}
