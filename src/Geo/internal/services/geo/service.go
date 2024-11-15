package geo

import (
	"context"
	"sort"
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/server"

	"github.com/mmcloughlin/geohash"
)

var _ server.Service = (*geoService)(nil)

const precision uint = 5

type GeoStorage interface {
	GetProfilesByGeohash(geohash string, gender model.Gender) []model.Profile
}

type ReactionServiceClient interface {
	GetReactedProfiles(profile_id int64) []model.Profile
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

	hash := geohash.EncodeWithPrecision(lat, lng, precision)
	geoProfiles := g.storage.GetProfilesByGeohash(hash, gender)
	reactedProfiles := g.reactionClient.GetReactedProfiles(profile_id)

	sort.Sort(model.ProfilesSortable(geoProfiles))
	sort.Sort(model.ProfilesSortable(reactedProfiles))

	result := g.findIntersection(geoProfiles, reactedProfiles)
	return result
}

func (g geoService) ChangeLocation(ctx context.Context, profile_id int64, lat, lng float64) error {
	// save to mongo
	panic("not implemented")
}

// find intersection by hash algorithm
func (f geoService) findIntersection(profiles1 []model.Profile, profiles2 []model.Profile) []model.Profile {
	m := make(map[int64]struct{}, len(profiles1))

	result := make([]model.Profile, 0, len(profiles1)/2)

	for _, v := range profiles1 {
		m[v.ID] = struct{}{}
	}

	for _, v := range profiles2 {
		_, ok := m[v.ID]
		if ok {
			result = append(result, v)
		}
	}

	return result
}
