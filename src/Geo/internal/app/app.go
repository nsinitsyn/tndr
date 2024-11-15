package app

import (
	"log"
	"log/slog"
	"os"
	"tinder-geo/internal/app/setup"
	"tinder-geo/internal/config"
	"tinder-geo/internal/infrastructure/clients"
	"tinder-geo/internal/infrastructure/database"
	"tinder-geo/internal/infrastructure/server"
	"tinder-geo/internal/services"
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

	storage := database.NewGeoStorage()
	reactionClient := clients.NewReactionServiceClient()
	service := services.NewGeoService(storage, reactionClient)

	if err := runGRPCServer(&config.GRPC, logger, service); err != nil {
		log.Fatal(err)
	}

	return nil
}

func runGRPCServer(config *config.GRPCConfig, logger *slog.Logger, service server.Service) error {
	serv := setup.NewGRPCServer(config, logger, service)
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
