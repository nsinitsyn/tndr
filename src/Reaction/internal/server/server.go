package server

import (
	"context"
	"log/slog"
	"strconv"
	"tinder-reaction/api/tinderpbv1"
	"tinder-reaction/internal/domain/model"

	"google.golang.org/grpc"
)

type ContextKey string

const ProfileIdCtxKey ContextKey = "profileId"
const GenderCtxKey ContextKey = "gender"

type Service interface {
	Like(ctx context.Context, profileId int64, gender model.Gender, likedProfileId int64) error
	Dislike(ctx context.Context, profileId int64, gender model.Gender, dislikedProfileId int64) error
}

type reactionServer struct {
	tinderpbv1.UnimplementedReactionServiceServer
	service Service
}

func Register(grpcServer *grpc.Server, service Service) {
	tinderpbv1.RegisterReactionServiceServer(grpcServer, &reactionServer{service: service})
}

func Stop(grpcServer *grpc.Server, logger *slog.Logger) {
	logger.Info("Stopping GRPC server...")
	grpcServer.GracefulStop()
}

// grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxIiwiR2VuZGVyIjoiTSIsImV4cCI6MTc2MzIwNzQxMywiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.VAVP65lIUhabxR4UknvQkRKiVCfu116cf3tZC8-dsfw' -plaintext -d '{"profile_id":456}' 172.24.48.1:2343 tinder.ReactionService/Like
func (s reactionServer) Like(ctx context.Context, req *tinderpbv1.LikeRequest) (*tinderpbv1.LikeResponse, error) {
	profileId, gender, err := parseArgumentsFromContext(ctx)
	if err != nil {
		return &tinderpbv1.LikeResponse{}, err
	}

	err = s.service.Like(ctx, profileId, gender, req.ProfileId)
	if err != nil {
		return &tinderpbv1.LikeResponse{}, err
	}

	return &tinderpbv1.LikeResponse{}, nil
}

func (s reactionServer) Dislike(ctx context.Context, req *tinderpbv1.DislikeRequest) (*tinderpbv1.DislikeResponse, error) {
	profileId, gender, err := parseArgumentsFromContext(ctx)
	if err != nil {
		return &tinderpbv1.DislikeResponse{}, err
	}

	err = s.service.Dislike(ctx, profileId, gender, req.ProfileId)
	return &tinderpbv1.DislikeResponse{}, err
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
