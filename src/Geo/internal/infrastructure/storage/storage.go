package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"tinder-geo/internal/config"
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/service"

	"github.com/redis/go-redis/v9"
)

const MAX_RETRIES = 10

var _ service.GeoStorage = (*geoStorage)(nil)

type geoStorage struct {
	config config.StorageConfig
}

func NewGeoStorage(config config.StorageConfig) geoStorage {
	return geoStorage{config: config}
}

func (s geoStorage) GetProfilesByGeohash(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     s.config.Addr,
		Password: s.config.Password,
		DB:       s.config.DB,
	})
	defer client.Close()

	profilesMap, err := client.HGetAll(ctx, fmt.Sprintf("geohash:%s:%s", geohash, gender.String())).Result()
	if err != nil {
		return nil, err
	}

	profiles := make([]model.Profile, 0, len(profilesMap))
	for _, p := range profilesMap {
		profile := model.Profile{}
		err = json.Unmarshal([]byte(p), &profile)

		if err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (s geoStorage) UpdateGeohash(ctx context.Context, profileId int64, gender model.Gender, geohash string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     s.config.Addr,
		Password: s.config.Password,
		DB:       s.config.DB,
	})
	defer client.Close()

	profileIdStr := strconv.FormatInt(profileId, 10)
	curGeo, _ := client.HGet(ctx, "profiles:geo", profileIdStr).Result()
	if geohash == curGeo {
		return nil
	}
	genderStr := gender.String()

	var versionKey = fmt.Sprintf("profile:version:%d", profileId)

	// Using CAS (check-and-set) for sync of concurrent UpdateGeohash and UpdateProfile calls
	for i := 0; i < MAX_RETRIES; i++ {
		err := client.Watch(ctx, func(tx *redis.Tx) error {
			n, err := tx.Get(ctx, versionKey).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			curGeo, err = client.HGet(ctx, "profiles:geo", profileIdStr).Result()
			if err != nil && err != redis.Nil {
				return err
			}

			// If first UpdateGeohash call before first UpdateProfile call then ignore it.
			// Of course we could add empty profile here and geohash for that.
			// But in this case other users can see empty profile in this geolocation. We would have to add a marker of an empty profile and take it when searching. It's overhead.
			// We just will wait until profile will be sync and update geohash later.
			if err == redis.Nil {
				// todo: как-то надо чекать ситуацию, если профиль по какой-то причине не дошел и уже не дойдет в GeoService. Тогда профиль никогда не будет показан другим пользователям!
				return nil
			}

			curGeoKey := fmt.Sprintf("geohash:%s:%s", curGeo, genderStr)
			newGeoKey := fmt.Sprintf("geohash:%s:%s", geohash, genderStr)
			// Get profile from cur:
			profile, _ := client.HGet(ctx, curGeoKey, profileIdStr).Result()

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				// Add profile to new geohash
				pipe.HSet(ctx, newGeoKey, profileIdStr, profile)
				// Remove profile from cur
				pipe.HDel(ctx, curGeoKey, profileIdStr)
				// Change index
				pipe.HSet(ctx, "profiles:geo", profileIdStr, geohash)

				// Update version
				pipe.Set(ctx, versionKey, strconv.FormatInt(n+1, 10), 0)
				return nil
			})
			return err
		}, versionKey)

		if err == redis.TxFailedErr {
			continue
		}
		if err == nil {
			return nil
		}
		return err
	}

	return errors.New("cannot update geohash. The maximum number of attempts to perform a transaction has been exceeded")
}

func (s geoStorage) UpdateProfile(ctx context.Context, gender model.Gender, profile model.Profile) error {
	client := redis.NewClient(&redis.Options{
		Addr:     s.config.Addr,
		Password: s.config.Password,
		DB:       s.config.DB,
	})
	defer client.Close()

	profileIdStr := strconv.FormatInt(profile.ID, 10)
	genderStr := gender.String()
	var versionKey = fmt.Sprintf("profile:version:%d", profile.ID)

	for i := 0; i < MAX_RETRIES; i++ {
		err := client.Watch(ctx, func(tx *redis.Tx) error {
			n, err := tx.Get(ctx, versionKey).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			curGeo, err := client.HGet(ctx, "profiles:geo", profileIdStr).Result()
			if err != nil && err != redis.Nil {
				return err
			}

			curGeoKey := fmt.Sprintf("geohash:%s:%s", curGeo, genderStr)
			if err == redis.Nil {
				// null - is special geohash for store new profiles until geolocation will be received
				curGeoKey = fmt.Sprintf("geohash:null:%s", genderStr)
			}

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				if err == redis.Nil {
					// if it's new profile then will be created index
					pipe.HSet(ctx, "profiles:geo", profileIdStr, "null")
				}

				_ = pipe.HSet(ctx, curGeoKey, profileIdStr, profile).Err()
				pipe.Set(ctx, versionKey, strconv.FormatInt(n+1, 10), 0)
				return nil
			})
			return err
		}, versionKey)

		if err == redis.TxFailedErr {
			continue
		}
		if err == nil {
			return nil
		}
		return err
	}

	return errors.New("cannot update profile. The maximum number of attempts to perform a transaction has been exceeded")
}
