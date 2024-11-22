package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"tinder-geo/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HTTPServer struct {
	srv    *http.Server
	config config.HTTPConfig
}

func NewHTTPServer(config config.HTTPConfig, promRegistry *prometheus.Registry) HTTPServer {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", config.Port)}
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(
		promRegistry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	srv.Handler = m

	return HTTPServer{
		srv:    srv,
		config: config,
	}
}

func (s HTTPServer) Run() error {
	slog.Info("HTTP server is running", slog.Int("port", s.config.Port))
	slog.Info("metrics are available via", slog.String("url", fmt.Sprintf("%s/metrics", s.srv.Addr)))

	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s HTTPServer) GracefulStop(ctx context.Context) {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "HTTP server shutdown error", slog.Any("error", err))
	}
}
