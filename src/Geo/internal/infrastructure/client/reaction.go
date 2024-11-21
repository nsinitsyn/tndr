package client

import (
	"context"
	"tinder-geo/internal/service"

	trace_utils "tinder-geo/internal/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var _ service.ReactionServiceClient = (*reactionServiceClient)(nil)

type reactionServiceClient struct {
}

func NewReactionServiceClient() reactionServiceClient {
	return reactionServiceClient{}
}

func (r reactionServiceClient) GetReactedProfiles(ctx context.Context, profileId int64) ([]int64, error) {
	tracer := otel.Tracer("reaction client")
	ctx, span := trace_utils.StartNewSpanWithCurrentFunctionName(
		ctx,
		tracer,
		attribute.Int64("tndr.params.profileId", profileId),
	)
	defer span.End()

	// not implemented yet - return empty slice
	result := []int64{}

	trace_utils.AddAttributesToCurrentSpan(ctx, attribute.Int64Slice("tndr.result.profileIds", result))

	return result, nil
}
