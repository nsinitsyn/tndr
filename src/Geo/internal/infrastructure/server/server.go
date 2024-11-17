package server

import (
	"context"
	"log/slog"
	"strconv"
	"tinder-geo/api/tinderpbv1"
	"tinder-geo/internal/domain/model"

	"google.golang.org/grpc"
)

type ContextKey string

const ProfileIdCtxKey ContextKey = "profileId"
const GenderCtxKey ContextKey = "gender"

type Service interface {
	GetProfilesByLocation(ctx context.Context, profile_id int64, gender model.Gender, lat, lng float64) []model.Profile // todo: model or dto?
	ChangeLocation(ctx context.Context, profile_id int64, lat, lng float64) error
}

type geoServer struct {
	tinderpbv1.UnimplementedGeoServiceServer
	service Service
}

func Register(gRPCServer *grpc.Server, service Service) {
	tinderpbv1.RegisterGeoServiceServer(gRPCServer, &geoServer{service: service})
}

func Stop(gRPCServer *grpc.Server, logger *slog.Logger) {
	logger.Info("Stopping GRPC server...")
	gRPCServer.GracefulStop()
}

// grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxIiwiR2VuZGVyIjoiTSIsImV4cCI6MTc2MzIwNzQxMywiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.VAVP65lIUhabxR4UknvQkRKiVCfu116cf3tZC8-dsfw' -plaintext -d '{"latitude":55.481, "longitude":37.288}' 172.24.48.1:2342 tinder.GeoService/GetProfilesByLocation
func (s *geoServer) GetProfilesByLocation(ctx context.Context, req *tinderpbv1.GetProfilesByLocationRequest) (*tinderpbv1.GetProfilesByLocationResponse, error) {
	profileId, gender, err := parseArgumentsFromContext(ctx)
	if err != nil {
		return &tinderpbv1.GetProfilesByLocationResponse{}, err
	}

	profiles := s.service.GetProfilesByLocation(ctx, profileId, gender, req.Latitude, req.Longitude)
	// todo: map profiles
	_ = profiles
	return &tinderpbv1.GetProfilesByLocationResponse{Profiles: nil}, nil
}

func (s *geoServer) ChangeLocation(ctx context.Context, req *tinderpbv1.ChangeLocationRequest) (*tinderpbv1.ChangeLocationResponse, error) {
	return &tinderpbv1.ChangeLocationResponse{}, nil
}

func parseArgumentsFromContext(ctx context.Context) (profileId int64, gender model.Gender, err error) {
	genderStr := ctx.Value(GenderCtxKey).(string)

	profileIdStr := ctx.Value(ProfileIdCtxKey).(string)
	profileId, err = strconv.ParseInt(profileIdStr, 10, 64)

	gender = model.M
	if genderStr == "F" {
		gender = model.F
	}

	return
}
