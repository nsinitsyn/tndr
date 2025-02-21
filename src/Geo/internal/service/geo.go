package service

import (
	"context"
	"fmt"
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/infrastructure/messaging"
	"tinder-geo/internal/server"

	trace_utils "tinder-geo/internal/trace"

	"github.com/mmcloughlin/geohash"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var _ server.Service = (*geoService)(nil)
var _ messaging.Service = (*geoService)(nil)

const PRECISION uint = 5

type GeoStorage interface {
	GetProfilesByGeohash(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error)
	UpdateGeohash(ctx context.Context, profileId int64, gender model.Gender, geohash string) error
	UpdateProfile(ctx context.Context, gender model.Gender, profile model.Profile) error
}

type ReactionServiceClient interface {
	GetReactedProfiles(ctx context.Context, profileId int64) ([]int64, error)
}

type geoService struct {
	storage        GeoStorage
	reactionClient ReactionServiceClient
}

func NewGeoService(storage GeoStorage, reactionClient ReactionServiceClient) geoService {
	return geoService{storage: storage, reactionClient: reactionClient}
}

func (g geoService) GetProfilesByLocation(ctx context.Context, profileId int64, gender model.Gender, lat, lng float64) ([]model.Profile, error) {
	ctx, span := trace_utils.StartNewSpanWithCurrentFunctionName(
		ctx,
		otel.Tracer("domain"),
		attribute.Float64("tndr.params.latitude", lat),
		attribute.Float64("tndr.params.longitude", lng),
	)
	defer span.End()

	geohash := geohash.EncodeWithPrecision(lat, lng, PRECISION)
	geoProfiles, err := g.storage.GetProfilesByGeohash(ctx, geohash, gender.Invert())
	if err != nil {
		return nil, err
	}

	excludedProfilesIds, err := g.reactionClient.GetReactedProfiles(ctx, profileId)
	if err != nil {
		return nil, err
	}

	excludedProfilesIds = append(excludedProfilesIds, profileId)
	result := g.excludeProfiles(geoProfiles, excludedProfilesIds)

	trace_utils.AddAttributesToCurrentSpan(ctx, attribute.String("tndr.result.profiles", fmt.Sprint(result)))

	return result, nil
}

func (g geoService) ChangeLocation(ctx context.Context, profileId int64, gender model.Gender, lat, lng float64) error {
	ctx, span := trace_utils.StartNewSpanWithCurrentFunctionName(
		ctx,
		otel.Tracer("domain"),
		attribute.Float64("tndr.params.latitude", lat),
		attribute.Float64("tndr.params.longitude", lng),
	)
	defer span.End()

	geohash := geohash.EncodeWithPrecision(lat, lng, PRECISION)
	err := g.storage.UpdateGeohash(ctx, profileId, gender, geohash)
	return err
}

func (g geoService) UpdateProfile(ctx context.Context, gender model.Gender, profile model.Profile) error {
	ctx, span := trace_utils.StartNewSpanWithCurrentFunctionName(
		ctx,
		otel.Tracer("domain"),
		attribute.String("tndr.params.profile", fmt.Sprint(profile)),
	)
	defer span.End()

	return g.storage.UpdateProfile(ctx, gender, profile)
}

func (f geoService) excludeProfiles(profiles []model.Profile, excludedProfilesIds []int64) []model.Profile {
	m := make(map[int64]struct{}, len(profiles))

	result := make([]model.Profile, 0, len(profiles)/2)

	for _, id := range excludedProfilesIds {
		m[id] = struct{}{}
	}

	for _, v := range profiles {
		_, ok := m[v.ID]
		if !ok {
			result = append(result, v)
		}
	}

	return result
}
