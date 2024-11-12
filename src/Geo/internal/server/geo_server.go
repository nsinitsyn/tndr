package server

import (
	"context"
	"tinder-geo/api/tinderpbv1"

	"google.golang.org/grpc"
)

type geoServer struct {
	tinderpbv1.UnimplementedGeoServiceServer
}

func Register(gRPCServer *grpc.Server) {
	tinderpbv1.RegisterGeoServiceServer(gRPCServer, &geoServer{})
}

func (s *geoServer) GetFeedByLocation(ctx context.Context, in *tinderpbv1.GetFeedByLocationRequest) (*tinderpbv1.GetFeedByLocationResponse, error) {
	return &tinderpbv1.GetFeedByLocationResponse{}, nil
}
