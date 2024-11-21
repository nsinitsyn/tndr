package main

import (
	"os"
	"os/signal"
	"syscall"
	"tinder-geo/internal/app"
)

func main() {
	closer := app.Run()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	closer()
}
