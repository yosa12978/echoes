package data

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
	"github.com/yosa12978/echoes/configs"
)

var (
	db     *sql.DB
	pgOnce sync.Once
)

func Postgres() *sql.DB {
	pgOnce.Do(func() {
		cfg := configs.Get()
		conn, err := sql.Open("postgres", cfg.Postgres)
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
