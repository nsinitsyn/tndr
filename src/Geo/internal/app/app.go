package app

import (
	"context"
	"log/slog"
	"os"
	"time"
	"tinder-geo/internal/config"
	"tinder-geo/internal/infrastructure/client"
	"tinder-geo/internal/infrastructure/messaging"
	"tinder-geo/internal/infrastructure/storage"
	"tinder-geo/internal/infrastructure/transport"
	"tinder-geo/internal/service"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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

	storage := storage.NewGeoStorage(config.Storage)
	reactionServiceClient := client.NewReactionServiceClient()
	service := service.NewGeoService(storage, reactionServiceClient)
	consumer := messaging.NewConsumer(config.Messaging, logger, storage)

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	grpcServer := transport.NewGRPCServer(config.GRPC, logger, service, promRegistry)
	httpServer := transport.NewHTTPServer(config.HTTP, logger, promRegistry)

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
		if err := grpcServer.Run(); err != nil {
			logger.Error("GRPC server starting error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	go func() {
		if err := httpServer.Run(); err != nil {
			logger.Error("HTTP server starting error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	return func() {
		cancel()

		// todo: smelly code for shutdown
		serverShutdown := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(serverShutdown)
		}()

		ctx, cancel = context.WithTimeout(context.Background(), GRACEFUL_SHUTDOWN_TIMEOUT_SEC*time.Second)
		defer cancel()
		httpServerShutdown := make(chan struct{})
		go func() {
			if err := httpServer.GracefulStop(ctx); err != nil {
				logger.Error("HTTP server shutdown error", slog.Any("error", err))
			}
			close(httpServerShutdown)
		}()

		timeout := time.After(GRACEFUL_SHUTDOWN_TIMEOUT_SEC * time.Second)
		stopped := 0
		for stopped != 3 {
			select {
			case <-consumingShutdown:
				consumingShutdown = nil
				stopped++
			case <-serverShutdown:
				serverShutdown = nil
				stopped++
			case <-httpServerShutdown:
				httpServerShutdown = nil
				stopped++
			case <-timeout:
				stopped = 3
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
