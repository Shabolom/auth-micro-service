package di

import (
	"auth-micro-service/internal/inmemory"
	"context"
	"fmt"
	"time"

	"auth-micro-service/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type DI struct {
	config *config.Config
	logger *zap.Logger

	inmemoryStorage *inmemory.SessionStorage

	ctx           context.Context
	pgConn        *pgxpool.Pool
	kafkaProducer *kafka.Writer
}

func New(ctx context.Context) *DI {
	return &DI{
		ctx: ctx,
	}
}

func (d *DI) GetInMemoryStorage() *inmemory.SessionStorage {
	if d.inmemoryStorage != nil {
		return d.inmemoryStorage
	}

	d.inmemoryStorage = inmemory.NewSessionStorage()
	d.inmemoryStorage.StartCleaner(d.ctx, time.Minute*5)
	return d.inmemoryStorage
}

func (d *DI) Config() *config.Config {
	if d.config != nil {
		return d.config
	}

	cfg, err := config.FromEnv()
	if err != nil {
		panic(fmt.Errorf("config from env: %w", err))
	}

	d.config = cfg
	return d.config
}

func (d *DI) Logger() *zap.Logger {
	if d.logger != nil {
		return d.logger
	}

	var logger *zap.Logger
	var err error

	if d.Config().Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(fmt.Errorf("create logger: %w", err))
	}

	logger = logger.With(
		zap.String("service", d.Config().ServiceName),
		zap.Bool("debug", d.Config().Debug),
	)

	d.logger = logger
	_ = zap.ReplaceGlobals(logger)

	return d.logger
}
