package clients

import (
	"context"
	"tinder-geo/internal/services"
)

var _ services.ReactionServiceClient = (*reactionServiceClient)(nil)

type reactionServiceClient struct {
}

func NewReactionServiceClient() *reactionServiceClient {
	return &reactionServiceClient{}
}

func (r reactionServiceClient) GetReactedProfiles(ctx context.Context, profileId int64) ([]int64, error) {
	// not implemented yet - return empty slice
	return []int64{}, nil
}
