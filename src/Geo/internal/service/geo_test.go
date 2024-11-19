package service

import (
	"context"
	"errors"
	"testing"
	"tinder-geo/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StorageMock struct {
	getProfilesByGeohash func(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error)
}

func (s StorageMock) GetProfilesByGeohash(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error) {
	return s.getProfilesByGeohash(ctx, geohash, gender)
}

func (s StorageMock) UpdateGeohash(ctx context.Context, profileId int64, gender model.Gender, geohash string) error {
	return nil
}

func (s StorageMock) UpdateProfile(ctx context.Context, gender model.Gender, profile model.Profile) error {
	return nil
}

type ReactionServiceClientMock struct {
	getReactedProfiles func(ctx context.Context, profileId int64) ([]int64, error)
}

func (r ReactionServiceClientMock) GetReactedProfiles(ctx context.Context, profileId int64) ([]int64, error) {
	return r.getReactedProfiles(ctx, profileId)
}

// go test -v -count=1 ./...

func TestGetProfiles(t *testing.T) {
	cases := map[string]struct {
		storageProfiles []model.Profile
		clientProfiles  []int64
		userID          int64
		storageError    error
		clientError     error
		result          []model.Profile
	}{
		"happy path": {
			storageProfiles: []model.Profile{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 5}, {ID: 4}, {ID: 8}},
			clientProfiles:  []int64{2, 5, 3},
			userID:          4,
			result:          []model.Profile{{ID: 1}, {ID: 8}},
		},
		"empty client": {
			storageProfiles: []model.Profile{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 5}, {ID: 4}, {ID: 8}},
			clientProfiles:  []int64{},
			userID:          1,
			result:          []model.Profile{{ID: 2}, {ID: 3}, {ID: 5}, {ID: 4}, {ID: 8}},
		},
		"client long": {
			storageProfiles: []model.Profile{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 5}, {ID: 4}, {ID: 8}},
			clientProfiles:  []int64{2, 5, 3, 1, 4, 9, 7},
			userID:          8,
			result:          []model.Profile{},
		},
		"storage error": {
			storageError:    errors.New("storage saving error"),
			storageProfiles: nil,
			clientProfiles:  []int64{2, 5, 3, 1, 4, 9, 7},
			userID:          8,
			result:          []model.Profile{},
		},
		"client error": {
			storageProfiles: []model.Profile{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 5}, {ID: 4}, {ID: 8}},
			clientError:     errors.New("client request error"),
			clientProfiles:  nil,
			userID:          8,
			result:          []model.Profile{},
		},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			storage := StorageMock{getProfilesByGeohash: func(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error) {
				return testCase.storageProfiles, testCase.storageError
			}}
			client := ReactionServiceClientMock{getReactedProfiles: func(ctx context.Context, profileId int64) ([]int64, error) {
				return testCase.clientProfiles, testCase.clientError
			}}
			service := NewGeoService(storage, client)
			actual, err := service.GetProfilesByLocation(context.Background(), testCase.userID, model.M, -1, -1)

			if testCase.storageError == nil && testCase.clientError == nil {
				require.NoError(t, err)
				assert.ElementsMatch(t, testCase.result, actual)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
