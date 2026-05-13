package di

import (
	"auth-micro-service/internal/redis-storage"
)

func (d *DI) GetRedisHandlers() *redisStorage.Redis {
	return redisStorage.NewRedisPublisher(d.NewRedisClient(), d.Logger())
}
