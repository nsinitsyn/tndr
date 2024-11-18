package clients

import (
	"tinder-geo/internal/services"
)

var _ services.ReactionServiceClient = (*reactionServiceClient)(nil)

type reactionServiceClient struct {
}

func NewReactionServiceClient() *reactionServiceClient {
	return &reactionServiceClient{}
}

func (r reactionServiceClient) GetReactedProfiles(profile_id int64) []int64 {
	// not implemented yet - return empty slice
	return []int64{}
}
