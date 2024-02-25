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
			Addr:       "0.0.0.0:8080",
			Postgres:   "postgres://root:root@localhost:5432/echoesdb?sslmode=disable",
			SessionKey: "34b6e28d-1b15-4739-9de0-8955586e56c2",
		}
		if addr, ok := os.LookupEnv("ADDR"); ok {
			config.Addr = addr
		}
		if cs, ok := os.LookupEnv("POSTGRES"); ok {
			config.Postgres = cs
		}
		if sessionKey, ok := os.LookupEnv("SESSION_KEY"); ok {
			config.SessionKey = sessionKey
		}
	})
	return *config
}
