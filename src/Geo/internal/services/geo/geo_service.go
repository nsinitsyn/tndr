package geo

import (
	"context"
	"sort"
	"tinder-geo/internal/domain/model"

	"github.com/mmcloughlin/geohash"
)

const precision uint = 5

type GetStorage interface {
	GetProfilesByGeohash(geohash string, sex bool) []model.Profile
}

type ReactionServiceClient interface {
	GetReactedProfiles(profile_id int64) []model.Profile
}

type feedService struct {
	storage               GetStorage
	reactionServiceClient ReactionServiceClient
}

func (f feedService) GetFeed(ctx context.Context, profile_id int64, lat, lng float64) []model.Profile {
	// todo:
	sex := true

	// todo: use ctx

	hash := geohash.EncodeWithPrecision(lat, lng, precision)
	geoProfiles := f.storage.GetProfilesByGeohash(hash, sex)
	reactedProfiles := f.reactionServiceClient.GetReactedProfiles(profile_id)

	sort.Sort(model.ProfilesSortable(geoProfiles))
	sort.Sort(model.ProfilesSortable(reactedProfiles))

	result := f.findIntersection(geoProfiles, reactedProfiles)
	return result
}

// find intersection by hash algorithm
func (f feedService) findIntersection(profiles1 []model.Profile, profiles2 []model.Profile) []model.Profile {
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
