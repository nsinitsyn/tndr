package service

import (
	"context"
	"tinder-reaction/internal/domain/model"
	"tinder-reaction/internal/infrastructure/transport"
	"tinder-reaction/internal/server"
)

var _ server.Service = (*reactionService)(nil)
var _ transport.Service = (*reactionService)(nil)

const PRECISION uint = 5

type ReactionStorage interface {
}

type reactionService struct {
	storage ReactionStorage
}

func NewReactionService(storage ReactionStorage) reactionService {
	return reactionService{storage: storage}
}

func (g reactionService) Like(ctx context.Context, profileId int64, gender model.Gender, likedProfileId int64) error {
	return nil
}

func (g reactionService) Dislike(ctx context.Context, profileId int64, gender model.Gender, dislikedProfileId int64) error {
	return nil
}

func (g reactionService) GetReactionsForProfile(ctx context.Context, profileId int64) ([]int64, error) {
	return []int64{100, 200, 300}, nil
}
