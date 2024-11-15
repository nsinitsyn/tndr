package app

import (
	"log"
	"log/slog"
	"os"
	"tinder-geo/internal/app/setup"
	"tinder-geo/internal/config"
	"tinder-geo/internal/infrastructure/clients"
	"tinder-geo/internal/infrastructure/database"
	"tinder-geo/internal/server"
	"tinder-geo/internal/services/geo"
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
	geoService := geo.NewGeoService(storage, reactionClient)

	if err := runGRPCServer(&config.GRPC, logger, geoService); err != nil {
		log.Fatal(err)
	}

	return nil
}

// grpcurl -plaintext 172.24.48.1:2342 list tinder.GeoService
// grpcurl -plaintext -d '{"geo_zone_id":5}' 172.24.48.1:2342 tinder.GeoService/GetFeedByLocation
func runGRPCServer(config *config.GRPCConfig, logger *slog.Logger, geoService server.Service) error {
	serv := setup.NewGRPCServer(config, logger, geoService)
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
