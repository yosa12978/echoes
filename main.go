package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yosa12978/echoes/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		log.Println(err.Error())
	}
}
