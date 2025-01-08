package app

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
	"tinder-reaction/internal/config"
	"tinder-reaction/internal/infrastructure/storage"
	"tinder-reaction/internal/infrastructure/transport"
	"tinder-reaction/internal/service"
	trace_utils "tinder-reaction/internal/trace"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel/trace"
)

const GRACEFUL_SHUTDOWN_TIMEOUT_SEC = 10

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run(ctx context.Context) {
	config := config.GetConfig()

	setupLogger(config.Service.Env)
	slog.Info("start...")

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	storage := storage.NewReactionStorage(config.Storage, promRegistry)
	service := service.NewReactionService(storage)

	if config.Tracing.Enabled {
		err := trace_utils.InitTracer(config.Tracing, config.Service)
		if err != nil {
			slog.Error("init tracer error", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("tracer initialized")
	}

	grpcServer := transport.NewGRPCServer(config.GRPC, service, promRegistry, config.Tracing.Enabled)
	httpServer := transport.NewHTTPServer(config.HTTP, service, promRegistry)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		if err := grpcServer.Run(); err != nil {
			slog.Error("GRPC server starting error", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("GRPC server stopped")
		wg.Done()
	}()

	go func() {
		if err := httpServer.Run(); err != nil {
			slog.Error("HTTP server starting error", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("HTTP server stopped")
		wg.Done()
	}()

	<-ctx.Done()
	go grpcServer.GracefulStop()
	go httpServer.GracefulStop(context.Background())
	go func() {
		storage.Close()
		wg.Done()
	}()

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

	slog.Info("application stopped")
}

func setupLogger(env string) {
	var handler slog.Handler
	switch env {
	case envLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envDev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	handler = LogHandler{handler}
	slog.SetDefault(slog.New(handler))
}

type LogHandler struct {
	slog.Handler
}

func (h LogHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := trace.SpanContextFromContext(ctx)

	if spanCtx.HasTraceID() {
		r.Add("trace_id", slog.StringValue(spanCtx.TraceID().String()))
	}
	if spanCtx.HasSpanID() {
		r.Add("span_id", slog.StringValue(spanCtx.SpanID().String()))
	}

	return h.Handler.Handle(ctx, r)
}
