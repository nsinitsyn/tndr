package app

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
	"tinder-geo/internal/config"
	"tinder-geo/internal/infrastructure/client"
	"tinder-geo/internal/infrastructure/messaging"
	"tinder-geo/internal/infrastructure/storage"
	"tinder-geo/internal/infrastructure/transport"
	"tinder-geo/internal/service"
	"tinder-geo/internal/trace"

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

	logger := setupLogger(config.Service.Env)
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

	if config.Tracing.Enabled {
		err := trace.InitTracer(config.Tracing, config.Service)
		if err != nil {
			logger.Error("init tracer error", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("tracer initialized")
	}

	grpcServer := transport.NewGRPCServer(config.GRPC, logger, service, promRegistry, config.Tracing.Enabled)
	httpServer := transport.NewHTTPServer(config.HTTP, logger, promRegistry)

	ctx, cancel := context.WithCancel(context.Background())
	consumingStarted := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		if err := consumer.StartConsume(ctx, consumingStarted); err != nil {
			logger.Error("fatal error", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("consuming stopped")
		wg.Done()
	}()

	go func() {
		if err := grpcServer.Run(); err != nil {
			logger.Error("GRPC server starting error", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("GRPC server stopped")
		wg.Done()
	}()

	go func() {
		if err := httpServer.Run(); err != nil {
			logger.Error("HTTP server starting error", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("HTTP server stopped")
		wg.Done()
	}()

	// Wait consuming start
	select {
	case <-consumingStarted:
		break
	case <-time.After(CONSUMING_START_TIMEOUT_SEC * time.Second):
		logger.Error("consuming start timeout expired")
		os.Exit(1)
	}

	return func() {
		cancel()
		go grpcServer.GracefulStop()
		go httpServer.GracefulStop(context.Background())

		stopped := make(chan struct{})
		go func() {
			wg.Wait()
			close(stopped)
		}()

		select {
		case <-stopped:
		case <-time.After(GRACEFUL_SHUTDOWN_TIMEOUT_SEC * time.Second):
			break
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
