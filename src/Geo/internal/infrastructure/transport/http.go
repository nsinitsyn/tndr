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
	logger *slog.Logger
}

func NewHTTPServer(config config.HTTPConfig, logger *slog.Logger, promRegistry *prometheus.Registry) HTTPServer {
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
		logger: logger,
	}
}

func (s HTTPServer) Run() error {
	s.logger.Info("HTTP server is running", slog.Int("port", s.config.Port))
	s.logger.Info("metrics are available via", slog.String("url", fmt.Sprintf("%s/metrics", s.srv.Addr)))

	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s HTTPServer) GracefulStop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err == nil {
		s.logger.Info("HTTP server stopped")
	}
	return err
}
