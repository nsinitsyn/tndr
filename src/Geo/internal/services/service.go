package services

import (
	"context"
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/infrastructure/server"

	"github.com/mmcloughlin/geohash"
)

var _ server.Service = (*geoService)(nil)

const PRECISION uint = 5

type GeoStorage interface {
	GetProfilesByGeohash(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error)
	UpdateGeohash(ctx context.Context, profileId int64, geohash string) error
}

type ReactionServiceClient interface {
	GetReactedProfiles(profile_id int64) []int64
}

type geoService struct {
	storage        GeoStorage
	reactionClient ReactionServiceClient
}

func NewGeoService(storage GeoStorage, reactionClient ReactionServiceClient) *geoService {
	return &geoService{storage: storage, reactionClient: reactionClient}
}

func (g geoService) GetProfilesByLocation(ctx context.Context, profile_id int64, gender model.Gender, lat, lng float64) []model.Profile {
	// todo: use ctx

	geohash := geohash.EncodeWithPrecision(lat, lng, PRECISION)
	geoProfiles, _ := g.storage.GetProfilesByGeohash(ctx, geohash, gender) // todo: err
	reactedProfilesIds := g.reactionClient.GetReactedProfiles(profile_id)

	// todo: удалить также реализацию интерфейса sort у профиля
	// sort.Sort(model.ProfilesSortable(geoProfiles))
	// sort.Sort(model.ProfilesSortable(reactedProfiles))

	result := g.findIntersection(geoProfiles, reactedProfilesIds)
	return result
}

func (g geoService) ChangeLocation(ctx context.Context, profileId int64, lat, lng float64) error {
	geohash := geohash.EncodeWithPrecision(lat, lng, PRECISION)
	err := g.storage.UpdateGeohash(ctx, profileId, geohash)
	return err
}

// find intersection by hash algorithm
func (f geoService) findIntersection(profiles []model.Profile, profilesIds []int64) []model.Profile {
	m := make(map[int64]struct{}, len(profiles))

	result := make([]model.Profile, 0, len(profiles)/2)

	for _, id := range profilesIds {
		m[id] = struct{}{}
	}

	for _, v := range profiles {
		_, ok := m[v.ID]
		if ok {
			result = append(result, v)
		}
	}

	return result
}
