package storage

import (
	"log/slog"
	"strconv"
	"strings"
	"tinder-reaction/internal/config"
	"tinder-reaction/internal/service"

	aero "github.com/aerospike/aerospike-client-go/v7"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	SET_NAME   string = "likes"
	INDEX_NAME string = "likes_profile_idx"
)

type BinName string

const (
	PROFILE_ID       BinName = "profile"
	LIKED_PROFILE_ID BinName = "liked"
	IS_LIKE          BinName = "like"
)

var _ service.ReactionStorage = (*reactionStorage)(nil)

type reactionStorage struct {
	config config.StorageConfig
	client *aero.Client
}

func NewReactionStorage(config config.StorageConfig, promRegistry *prometheus.Registry) reactionStorage {
	client, err := aero.NewClient(config.Hostname, config.Port)
	if err != nil {
		panic(err)
	}
	// todo: call createIndex only once
	if err := createIndex(client, config.Namespace); err != nil {
		panic(err)
	}
	return reactionStorage{config: config, client: client}
}

func (s reactionStorage) Close() {
	s.client.Close()
}

func createIndex(client *aero.Client, namespace string) error {
	task, err := client.CreateIndex(nil, namespace, SET_NAME, INDEX_NAME, string(PROFILE_ID), aero.NUMERIC)
	if err != nil {
		return err
	}
	err = <-task.OnComplete()
	return err
}

func (s reactionStorage) AddLike(profileId int64, likedProfileId int64) error {
	return s.addReaction(profileId, likedProfileId, true)
}

func (s reactionStorage) AddDislike(profileId int64, likedProfileId int64) error {
	return s.addReaction(profileId, likedProfileId, false)
}

func (s reactionStorage) GetReactions(profileId int64) ([]int64, error) {
	stmt := aero.NewStatement(s.config.Namespace, SET_NAME, string(LIKED_PROFILE_ID))
	stmt.SetFilter(aero.NewEqualFilter(string(PROFILE_ID), profileId))
	recordSet, err := s.client.Query(nil, stmt)
	if err != nil {
		return nil, err
	}
	// todo: calc count of records
	result := make([]int64, 0, 100)
	for record := range recordSet.Results() {
		if record != nil {
			if record.Err != nil {
				slog.Error("error during reactions iterations from aerospike", slog.Any("error", record.Err))
			} else {
				result = append(result, int64(record.Record.Bins[string(LIKED_PROFILE_ID)].(int)))
			}
		}
	}

	return result, nil
}

func (s reactionStorage) addReaction(profileId, likedProfileId int64, like bool) error {
	key, err := aero.NewKey(s.config.Namespace, SET_NAME, createKey(profileId, likedProfileId))
	if err != nil {
		return err
	}
	bins := aero.BinMap{
		string(PROFILE_ID):       profileId,
		string(LIKED_PROFILE_ID): likedProfileId,
		string(IS_LIKE):          like,
	}
	err = s.client.Put(nil, key, bins)
	return err
}

func createKey(user1, user2 int64) string {
	user1Str := strconv.FormatInt(user1, 10)
	user2Str := strconv.FormatInt(user2, 10)

	var sb strings.Builder
	sb.Grow(len(user1Str) + len(user2Str) + 1)

	sb.WriteString(user1Str)
	sb.WriteString("|")
	sb.WriteString(user2Str)

	return sb.String()
}
