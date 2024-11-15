package clients

import (
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/services"
)

var _ services.ReactionServiceClient = (*reactionServiceClient)(nil)

type reactionServiceClient struct {
}

func NewReactionServiceClient() *reactionServiceClient {
	return &reactionServiceClient{}
}

func (r reactionServiceClient) GetReactedProfiles(profile_id int64) []model.Profile {
	return nil
}
