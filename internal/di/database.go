package di

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const initTimeout = 3 * time.Second

func (d *DI) GetPgDatabase() *pgxpool.Pool {
	if d.pgConn != nil {
		return d.pgConn
	}

	ctx, cancel := context.WithTimeout(context.Background(), initTimeout)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(d.Config().DatabaseURL())
	if err != nil {
		log.Fatal(fmt.Errorf("parse database url: %w", err))
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("create pg pool: %w", err))
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal(fmt.Errorf("ping database: %w", err))
	}

	d.pgConn = pool
	return d.pgConn
}
