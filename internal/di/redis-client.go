package di

import (
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (d *DI) NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     d.Config().RedisDSN(),
		Password: d.Config().Redis.RedisPassword,
		DB:       0,

		PoolSize:     4,
		MinIdleConns: 2,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		MaxRetries:  2,
		PoolTimeout: 2 * time.Second,
	})

	if err := rdb.Ping(d.ctx).Err(); err != nil {
		d.Logger().Fatal("redis ping err :", zap.Error(err))
	}

	return rdb
}
