package di

import "auth-micro-service/internal/redis"

func (d *DI) GetRedisHandlers() *redis.Redis {
	return redis.NewRedisPublisher(d.NewRedisClient(), d.Logger())
}
