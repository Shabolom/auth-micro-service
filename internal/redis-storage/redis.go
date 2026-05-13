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
	if r.client == nil {
		return nil
	}

	return r.client.Close()
}
