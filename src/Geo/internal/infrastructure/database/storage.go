package database

import (
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/services/geo"
)

var _ geo.GeoStorage = (*geoStorage)(nil)

type geoStorage struct {
}

func NewGeoStorage() *geoStorage {
	return &geoStorage{}
}

func (s geoStorage) GetProfilesByGeohash(geohash string, gender model.Gender) []model.Profile {
	return nil
}
