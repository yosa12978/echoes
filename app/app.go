package app

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/yosa12978/echoes/configs"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/session"
)

func init() {
	if os.Getenv("DEBUG") == "true" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}
}

func Run(ctx context.Context) error {
	logger := logging.New("app.Run")
	conn := data.Postgres()
	defer conn.Close()

	rdb := data.Redis(ctx)
	defer rdb.Close()

	session.SetupStore()

	cfg := configs.Get()

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: NewRouter(ctx),
	}

	errch := make(chan error, 1)
	go func() {
		logger.Printf("server listening @ %s", cfg.Addr)
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
