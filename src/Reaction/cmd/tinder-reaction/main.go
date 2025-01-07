package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tinder-reaction/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		stop()
	}()
	app.Run(ctx)
}
