package redisStorage

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (r *Redis) CheckSessionStatus(ctx context.Context, jti string) error {
	fmt.Println(11111111)
	result, err := r.client.HGet(ctx, jti, REVOKE).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Info("session status not exist", zap.String("jti", jti))
			return redis.Nil
		}

		return err
	}

	status, err := strconv.ParseBool(result)
	if err != nil {
		fmt.Println(2222222)
		r.logger.Info("session status not bool", zap.String("jti", jti))
		return err
	}

	if status {
		r.logger.Info("session status REVOKE true", zap.String("jti", jti))
		return errors.New("session revoked")
	}

	fmt.Println(4444444, status)
	return nil
}
