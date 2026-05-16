package di

import (
	"auth-micro-service/internal/redis-storage"
)

func (d *DI) GetRedisHandlers() *redisStorage.Redis {
	if d.redis != nil {
		return d.redis
	}

	return redisStorage.NewRedisPublisher(d.NewRedisClient(), d.Logger())
}
