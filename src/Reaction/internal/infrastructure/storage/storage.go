package storage

import (
	"tinder-reaction/internal/config"
	"tinder-reaction/internal/service"

	"github.com/prometheus/client_golang/prometheus"
)

var _ service.ReactionStorage = (*reactionStorage)(nil)

type reactionStorage struct {
	config config.StorageConfig
}

func NewReactionStorage(config config.StorageConfig, promRegistry *prometheus.Registry) reactionStorage {
	return reactionStorage{config: config}
}
