package configs

import (
	"os"
	"sync"

	"github.com/yosa12978/echoes/types"
)

var (
	config *types.Config
	once   sync.Once
)

func Get() types.Config {
	once.Do(func() {
		config = &types.Config{
			Addr:     "0.0.0.0:8080",
			Postgres: "postgres://root:root@localhost:5432/echoesdb?sslmode=disable",
		}
		if addr, ok := os.LookupEnv("ADDR"); ok {
			config.Addr = addr
		}
		if cs, ok := os.LookupEnv("POSTGRES"); ok {
			config.Postgres = cs
		}
	})
	return *config
}
