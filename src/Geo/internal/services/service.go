package services

import (
	"context"
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/server"

	"github.com/mmcloughlin/geohash"
)

var _ server.Service = (*geoService)(nil)

const PRECISION uint = 5

type GeoStorage interface {
	GetProfilesByGeohash(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error)
	UpdateGeohash(ctx context.Context, profileId int64, gender model.Gender, geohash string) error
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

func (g geoService) GetProfilesByLocation(ctx context.Context, profileId int64, gender model.Gender, lat, lng float64) []model.Profile {
	// todo: use ctx

	geohash := geohash.EncodeWithPrecision(lat, lng, PRECISION)
	geoProfiles, _ := g.storage.GetProfilesByGeohash(ctx, geohash, gender.Invert()) // todo: err
	excludedProfilesIds := g.reactionClient.GetReactedProfiles(profileId)

	// todo: удалить также реализацию интерфейса sort у профиля
	// sort.Sort(model.ProfilesSortable(geoProfiles))
	// sort.Sort(model.ProfilesSortable(reactedProfiles))

	excludedProfilesIds = append(excludedProfilesIds, profileId)
	result := g.excludeProfiles(geoProfiles, excludedProfilesIds)

	// todo: feature: if profiles wasn't found in specified geohash then to find in neighboring geohashes

	return result
}

func (g geoService) ChangeLocation(ctx context.Context, profileId int64, gender model.Gender, lat, lng float64) error {
	geohash := geohash.EncodeWithPrecision(lat, lng, PRECISION)
	err := g.storage.UpdateGeohash(ctx, profileId, gender, geohash)
	return err
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
