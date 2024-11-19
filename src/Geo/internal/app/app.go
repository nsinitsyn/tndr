package app

import (
	"context"
	"log/slog"
	"os"
	"time"
	"tinder-geo/internal/config"
	"tinder-geo/internal/infrastructure/clients"
	"tinder-geo/internal/infrastructure/messaging"
	"tinder-geo/internal/infrastructure/storage"
	"tinder-geo/internal/infrastructure/transport"
	"tinder-geo/internal/services"
)

const GRACEFUL_SHUTDOWN_TIMEOUT_SEC = 10
const CONSUMING_START_TIMEOUT_SEC = 15

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() (closer func()) {
	config := config.GetConfig()

	logger := setupLogger(config.Env)
	logger.Info("start...")

	storage := storage.NewGeoStorage(&config.Storage)
	reactionServiceClient := clients.NewReactionServiceClient()
	service := services.NewGeoService(storage, reactionServiceClient)
	consumer := messaging.NewConsumer(config.Messaging, logger, storage)
	server := transport.NewServer(&config.GRPC, logger, service)

	ctx, cancel := context.WithCancel(context.Background())
	consumingStarted := make(chan struct{})
	consumingShutdown := make(chan struct{})
	go func() {
		if err := consumer.StartConsume(ctx, consumingStarted); err != nil {
			logger.Error("fatal error", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("consuming stopped")
		close(consumingShutdown)
	}()

	select {
	case <-consumingStarted:
		break
	case <-time.After(CONSUMING_START_TIMEOUT_SEC * time.Second):
		logger.Error("fatal error: consuming start timeout expired")
		os.Exit(1)
	}

	go func() {
		if err := server.Run(); err != nil {
			logger.Error("fatal error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	return func() {
		cancel()

		serverShutdown := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(serverShutdown)
		}()

		timeout := time.After(GRACEFUL_SHUTDOWN_TIMEOUT_SEC * time.Second)
		stopped := 0
		for stopped != 2 {
			select {
			case <-consumingShutdown:
				consumingShutdown = nil
				stopped++
			case <-serverShutdown:
				serverShutdown = nil
				stopped++
			case <-timeout:
				stopped = 2
			}
		}

		logger.Info("application stopped")
	}
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
