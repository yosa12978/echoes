package data

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"github.com/yosa12978/echoes/config"
)

var (
	db     *sql.DB
	pgOnce sync.Once
)

func Postgres() *sql.DB {
	pgOnce.Do(func() {
		cfg := config.Get()
		s := fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=%s",
			cfg.Postgres.User,
			cfg.Postgres.Pass,
			cfg.Postgres.Addr,
			cfg.Postgres.DB,
			cfg.Postgres.SSLMode,
		)
		conn, err := sql.Open("postgres", s)
		if err != nil {
			panic(err)
		}
		if err := conn.Ping(); err != nil {
			panic(err)
		}
		db = conn
	})
	return db
}

type pgPinger struct {
	pg *sql.DB
}

func NewPgPinger() Pinger {
	return &pgPinger{
		pg: db,
	}
}

func (p *pgPinger) Ping(ctx context.Context) error {
	return p.pg.Ping()
}
