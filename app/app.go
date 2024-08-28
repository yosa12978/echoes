package app

import (
	"context"
	"os"
	"time"

	"github.com/yosa12978/echoes/config"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/session"
)

func Run(ctx context.Context) error {
	logger := logging.NewJsonLogger(os.Stdout)
	conn := data.Postgres()
	defer conn.Close()

	rdb := data.Redis(ctx)
	defer rdb.Close()

	session.SetupStore()

	cfg := config.Get()

	server := newServer(
		ctx,
		cfg.Server.Addr,
		logger,
	)

	errch := make(chan error, 1)
	go func() {
		logger.Info("server listening", "addr", cfg.Server.Addr)
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
