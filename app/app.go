package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/yosa12978/echoes/configs"
	"github.com/yosa12978/echoes/data"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func Run(ctx context.Context) error {
	conn := data.Postgres()
	defer conn.Close()

	cfg := configs.Get()

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: NewRouter(ctx),
	}

	errch := make(chan error, 1)
	go func() {
		log.Printf("server listening @ %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil {
			errch <- err
		}
		close(errch)
	}()

	var err error
	select {
	case err = <-errch:
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err = server.Shutdown(timeout)
	}
	return err
}
