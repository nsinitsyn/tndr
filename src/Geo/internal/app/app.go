package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinder-geo/internal/app/setup"
	"tinder-geo/internal/config"
	"tinder-geo/internal/infrastructure/clients"
	"tinder-geo/internal/infrastructure/messaging"
	"tinder-geo/internal/infrastructure/storage"
	"tinder-geo/internal/services"
)

const GRACEFUL_SHUTDOWN_TIMEOUT_SEC = 10

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() {
	config := config.GetConfig()

	logger := setupLogger(config.Env)
	_ = logger

	storage := storage.NewGeoStorage(&config.Storage)
	reactionServiceClient := clients.NewReactionServiceClient()
	service := services.NewGeoService(storage, reactionServiceClient)
	consumer := messaging.NewConsumer(config.Messaging, logger, storage)
	// todo: заменить setup!!!
	grpcServer := setup.NewGRPCServer(&config.GRPC, logger, service)

	ctx, cancel := context.WithCancel(context.Background())
	consumingShutdown := make(chan struct{})
	go func() {
		if err := consumer.StartConsume(ctx); err != nil {
			log.Fatal(err)
		}
		logger.Info("consuming messaging stopped")
		close(consumingShutdown)
	}()

	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	cancel()

	serverShutdown := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
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

	logger.Info("stopped")
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
