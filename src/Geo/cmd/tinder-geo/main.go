package main

import (
	"os"
	"os/signal"
	"syscall"
	"tinder-geo/internal/app"
)

// MAC OS - локальный запуск инфры и приложения:
// >cd cicd/debug
// >docker compose up -d
// >cd src/Geo
// >export CONFIG_PATH='config/config.yaml'; export CGO_ENABLED=1; go run cmd/tinder-geo/main.go

func main() {
	closer := app.Run()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	closer()
}
