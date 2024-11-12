package app

import (
	"log"
	"log/slog"
	"os"
	"tinder-geo/internal/app/setup"
	"tinder-geo/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() error {
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
	serv := setup.NewGRPCServer(config, logger)
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
