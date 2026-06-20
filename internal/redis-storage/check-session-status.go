package redisStorage

import (
	"auth-micro-service/pkg/shortcut"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (r *Redis) CheckSessionStatus(ctx context.Context, jti string) error {
	result, err := r.client.HGet(ctx, jti, REVOKE).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Info("session status not exist", zap.String("jti", jti))
			return shortcut.ErrSessionNotFound
		}

		return err
	}

	status, err := strconv.ParseBool(result)
	if err != nil {
		r.logger.Info("session status not bool", zap.String("jti", jti))
		return err
	}

	if status {
		r.logger.Info("session status REVOKE true", zap.String("jti", jti))
		return shortcut.ErrSessionRevoked
	}

	expiredUnixStr, err := r.client.HGet(ctx, jti, EXPIRE).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Info("session status not exist", zap.String("jti", jti))
			return redis.Nil
		}

		return err
	}

	expiredUnix, err := strconv.ParseInt(expiredUnixStr, 10, 64)
	if err != nil {
		r.logger.Info("session status not int", zap.String("jti", jti))
		return err
	}

	expiredTime := time.Unix(expiredUnix, 0)
	if time.Now().After(expiredTime) {
		r.logger.Info("session status expired true", zap.String("jti", jti))
		return shortcut.ErrSessionExpired
	}

	return nil
}
