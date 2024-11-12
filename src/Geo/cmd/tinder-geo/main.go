package main

import (
	"log"
	"log/slog"
	"os"
	"tinder-geo/internal/config"
	"tinder-geo/internal/server"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// $env:CONFIG_PATH = '././config/config.yaml'; go run cmd/tinder-geo/main.go

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	// code := geohash.EncodeWithPrecision(55.7893, 37.7717, 5)
	// fmt.Println(code)
	// // clean arh: app в коде Олега Козырева и видео art of development
	// // grpc server
	// // kafka client
	// // mongo client
	// // match service http client - 100 ошибок go - как закрывать body http клиентов правильно
	// // business logic - searching nearby profiles by geohash
	// // observability - 2 видео art of development
	// // graceful shutdown
	// // tests
}

func run() error {
	config := config.GetConfig()

	logger := setupLogger(config.Env)
	_ = logger

	if err := runGRPCServer(&config.GRPC, logger); err != nil {
		log.Fatal(err)
	}

	return nil
}

// grpcurl -plaintext 172.24.48.1:2342 list tinder.GeoService
// grpcurl -plaintext -d '{"geo_zone_id":5}' 172.24.48.1:2342 tinder.GeoService/GetFeedByLocation
func runGRPCServer(config *config.GRPCConfig, logger *slog.Logger) error {
	serv := server.NewGRPCServer(config, logger)
	if err := serv.Run(); err != nil {
		return err
	}

	return nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
