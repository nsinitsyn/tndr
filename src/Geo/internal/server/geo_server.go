package server

import (
	"context"
	"log/slog"
	"tinder-geo/api/tinderpbv1"
	"tinder-geo/internal/domain/model"

	"google.golang.org/grpc"
)

type Service interface {
	GetFeed(ctx context.Context, profile_id int64, lat, lng float64) []model.Profile // todo: model or dto?
}

type geoServer struct {
	tinderpbv1.UnimplementedGeoServiceServer
	service Service
}

func Register(gRPCServer *grpc.Server) {
	tinderpbv1.RegisterGeoServiceServer(gRPCServer, &geoServer{})
}

func Stop(gRPCServer *grpc.Server, logger *slog.Logger) {
	logger.Info("Stopping GRPC server...")
	gRPCServer.GracefulStop()
}

func (s *geoServer) GetFeedByLocation(ctx context.Context, req *tinderpbv1.GetFeedByLocationRequest) (*tinderpbv1.GetFeedByLocationResponse, error) {
	return &tinderpbv1.GetFeedByLocationResponse{}, nil
}

func (s *geoServer) ChangeLocation(ctx context.Context, req *tinderpbv1.ChangeLocationRequest) (*tinderpbv1.ChangeLocationResponse, error) {
	return &tinderpbv1.ChangeLocationResponse{}, nil
}
