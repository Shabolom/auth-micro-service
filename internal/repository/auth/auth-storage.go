package auth

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Storage {
	return &Storage{
		conn: conn,
	}
}
