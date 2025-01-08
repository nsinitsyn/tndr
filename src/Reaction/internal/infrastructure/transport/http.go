package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"tinder-reaction/internal/config"
	"tinder-reaction/internal/infrastructure/transport/model"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service interface {
	GetReactionsForProfile(ctx context.Context, profileId int64) ([]int64, error)
}

type HTTPServer struct {
	srv     *http.Server
	config  config.HTTPConfig
	service Service
}

func NewHTTPServer(
	config config.HTTPConfig,
	service Service,
	promRegistry *prometheus.Registry) HTTPServer {
	m := http.NewServeMux()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: m,
	}

	server := HTTPServer{
		srv:     srv,
		config:  config,
		service: service,
	}

	m.Handle("/metrics", promhttp.HandlerFor(
		promRegistry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	m.HandleFunc("GET /reactions/{profileId}", server.getReactionsHandler)

	return server
}

func (s HTTPServer) getReactionsHandler(w http.ResponseWriter, r *http.Request) {
	// todo: в middleware проверить jwt - должен быть Admin или Service:Geo
	// todo: здесь нужно получить traceId от geo сервиса и продолжить регистрировать спаны в нем
	profileIdParam := r.PathValue("profileId")
	profileId, err := strconv.ParseInt(profileIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := model.ErrorResponse{Error: "Parameter 'profileId' must be of type int64"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// slog.Info("got data", slog.Int64("profileId", profileId))
	profiles, err := s.service.GetReactionsForProfile(r.Context(), profileId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("get reaction error", slog.Any("error", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)
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
