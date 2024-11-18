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
	GetProfilesByLocation(ctx context.Context, profileId int64, gender model.Gender, lat, lng float64) []model.Profile // todo: model or dto?
	ChangeLocation(ctx context.Context, profileId int64, gender model.Gender, lat, lng float64) error
}

type geoServer struct {
	tinderpbv1.UnimplementedGeoServiceServer
	service Service
}

func Register(grpcServer *grpc.Server, service Service) {
	tinderpbv1.RegisterGeoServiceServer(grpcServer, &geoServer{service: service})
}

func Stop(grpcServer *grpc.Server, logger *slog.Logger) {
	logger.Info("Stopping GRPC server...")
	grpcServer.GracefulStop()
}

// grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxIiwiR2VuZGVyIjoiTSIsImV4cCI6MTc2MzIwNzQxMywiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.VAVP65lIUhabxR4UknvQkRKiVCfu116cf3tZC8-dsfw' -plaintext -d '{"latitude":55.481, "longitude":37.288}' 172.24.48.1:2342 tinder.GeoService/GetProfilesByLocation
func (s *geoServer) GetProfilesByLocation(ctx context.Context, req *tinderpbv1.GetProfilesByLocationRequest) (*tinderpbv1.GetProfilesByLocationResponse, error) {
	profileId, gender, err := parseArgumentsFromContext(ctx)
	if err != nil {
		return &tinderpbv1.GetProfilesByLocationResponse{}, err
	}

	profiles := s.service.GetProfilesByLocation(ctx, profileId, gender, req.Latitude, req.Longitude)
	profilesDtos := make([]*tinderpbv1.LocationProfileDto, 0, len(profiles))
	for _, v := range profiles {
		profilesDtos = append(profilesDtos, &tinderpbv1.LocationProfileDto{
			ProfileId:   v.ID,
			Name:        v.Name,
			Description: v.Description,
			PhotoUrls:   v.Photos,
		})
	}
	return &tinderpbv1.GetProfilesByLocationResponse{Profiles: profilesDtos}, nil
}

func (s *geoServer) ChangeLocation(ctx context.Context, req *tinderpbv1.ChangeLocationRequest) (*tinderpbv1.ChangeLocationResponse, error) {
	profileId, gender, err := parseArgumentsFromContext(ctx)
	if err != nil {
		return &tinderpbv1.ChangeLocationResponse{}, err
	}

	err = s.service.ChangeLocation(ctx, profileId, gender, req.Latitude, req.Longitude)
	return &tinderpbv1.ChangeLocationResponse{}, err
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
