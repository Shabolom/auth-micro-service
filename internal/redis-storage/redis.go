package redisStorage

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	USERID = "UserID"
	EXPIRE = "ExpiresAt"
	REVOKE = "Revoke"
)

type Redis struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisPublisher(client *redis.Client, logger *zap.Logger) *Redis {
	return &Redis{
		client: client,
		logger: logger,
	}
}

func (r *Redis) Close() error {
	err := r.client.Close()
	if err != nil {
		r.logger.Error("Error closing redis client", zap.Error(err))
		return err
	}

	return nil
}
