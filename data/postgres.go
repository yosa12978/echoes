package data

import (
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
