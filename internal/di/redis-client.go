package di

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (d *DI) NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     d.Config().RedisDSN(),
		Password: d.Config().Redis.RedisPassword,
		DB:       0,

		PoolSize:     10,
		MinIdleConns: 2,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		MaxRetries:  2,
		PoolTimeout: 2 * time.Second,

		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: time.Hour,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		d.Logger().Fatal("redis ping err :", zap.Error(err))
	}

	return rdb
}
