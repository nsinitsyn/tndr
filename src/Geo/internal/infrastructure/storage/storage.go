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

	"github.com/redis/go-redis/extra/redisprometheus/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	trace_utils "tinder-geo/internal/trace"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const MAX_RETRIES = 10

var _ service.GeoStorage = (*geoStorage)(nil)

type geoStorage struct {
	config config.StorageConfig
	client *redis.Client
}

func NewGeoStorage(config config.StorageConfig, promRegistry *prometheus.Registry) geoStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	collector := redisprometheus.NewCollector("geo", "redis", client)
	promRegistry.MustRegister(collector)

	if err := redisotel.InstrumentTracing(client); err != nil {
		panic(err)
	}

	return geoStorage{config: config, client: client}
}

func (s geoStorage) Close() error {
	return s.client.Close()
}

func (s geoStorage) GetProfilesByGeohash(ctx context.Context, geohash string, gender model.Gender) ([]model.Profile, error) {
	tracer := otel.Tracer("storage")
	ctx, span := trace_utils.StartNewSpanWithCurrentFunctionName(
		ctx,
		tracer,
		attribute.String("tndr.params.geohash", geohash),
		attribute.String("tndr.params.gender", gender.String()),
	)
	defer span.End()

	profilesMap, err := s.client.HGetAll(ctx, fmt.Sprintf("geohash:%s:%s", geohash, gender.String())).Result()
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

	trace_utils.AddAttributesToCurrentSpan(ctx, attribute.String("tndr.result.profiles", fmt.Sprint(profiles)))

	return profiles, nil
}

func (s geoStorage) UpdateGeohash(ctx context.Context, profileId int64, gender model.Gender, geohash string) error {
	ctx, span := trace_utils.StartNewSpanWithCurrentFunctionName(
		ctx,
		otel.Tracer("domain"),
		attribute.Int64("tndr.params.profileId", profileId),
		attribute.String("tndr.params.gender", gender.String()),
		attribute.String("tndr.params.geohash", geohash),
	)
	defer span.End()

	profileIdStr := strconv.FormatInt(profileId, 10)
	curGeo, _ := s.client.HGet(ctx, "profiles:geo", profileIdStr).Result()
	if geohash == curGeo {
		return nil
	}
	genderStr := gender.String()

	var versionKey = fmt.Sprintf("profile:version:%d", profileId)

	// Using CAS (check-and-set) for sync of concurrent UpdateGeohash and UpdateProfile calls
	for i := 0; i < MAX_RETRIES; i++ {
		err := s.client.Watch(ctx, func(tx *redis.Tx) error {
			n, err := tx.Get(ctx, versionKey).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			curGeo, err = s.client.HGet(ctx, "profiles:geo", profileIdStr).Result()
			if err != nil && err != redis.Nil {
				return err
			}

			// If first UpdateGeohash call before first UpdateProfile call then ignore it.
			// Of course we could add empty profile here and geohash for that.
			// But in this case other users can see empty profile in this geolocation. We would have to add a marker of an empty profile and take it when searching. It's overhead.
			// We just will wait until profile will be sync and update geohash later.
			if err == redis.Nil {
				return nil
			}

			curGeoKey := fmt.Sprintf("geohash:%s:%s", curGeo, genderStr)
			newGeoKey := fmt.Sprintf("geohash:%s:%s", geohash, genderStr)
			// Get profile from cur:
			profile, _ := s.client.HGet(ctx, curGeoKey, profileIdStr).Result()

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
	ctx, span := trace_utils.StartNewSpanWithCurrentFunctionName(
		ctx,
		otel.Tracer("domain"),
		attribute.String("tndr.params.gender", gender.String()),
		attribute.String("tndr.params.profile", fmt.Sprint(profile)),
	)
	defer span.End()

	profileIdStr := strconv.FormatInt(profile.ID, 10)
	genderStr := gender.String()
	var versionKey = fmt.Sprintf("profile:version:%d", profile.ID)

	for i := 0; i < MAX_RETRIES; i++ {
		err := s.client.Watch(ctx, func(tx *redis.Tx) error {
			n, err := tx.Get(ctx, versionKey).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			curGeo, err := s.client.HGet(ctx, "profiles:geo", profileIdStr).Result()
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
