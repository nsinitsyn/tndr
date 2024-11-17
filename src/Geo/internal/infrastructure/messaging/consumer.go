package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"
	"tinder-geo/internal/config"
	"tinder-geo/internal/domain/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const BATCH_SIZE int = 10
const READ_MESSAGE_TIMEOUT_SEC = 5

type GeoStorage interface {
	UpdateProfile(ctx context.Context, gender model.Gender, profile model.Profile) error
}

type kafkaConsumer struct {
	config  config.MessagingConfig
	logger  *slog.Logger
	storage GeoStorage
}

func NewConsumer(config config.MessagingConfig, logger *slog.Logger, storage GeoStorage) *kafkaConsumer {
	return &kafkaConsumer{config: config, logger: logger, storage: storage}
}

func (k kafkaConsumer) StartConsume(ctx context.Context) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  k.config.Servers,
		"group.id":           k.config.Group,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return err
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{k.config.Topic}, nil)

	k.logger.Info("start consuming...")

	batch := make([]ProfileDto, 0, BATCH_SIZE)

	for !errors.Is(ctx.Err(), context.Canceled) {
		msg, err := consumer.ReadMessage(READ_MESSAGE_TIMEOUT_SEC * time.Second)
		if err == nil {
			var profileDto ProfileDto
			err := json.Unmarshal(msg.Value, &profileDto)
			if err != nil {
				k.logger.Error(
					"error decoding message",
					slog.Any("error", err),
					slog.Any("message", string(msg.Value)))

				// todo: publish to dead letter queue...
				continue
			}

			batch = append(batch, profileDto)

			if len(batch) == BATCH_SIZE {
				k.processBatch(ctx, batch)
				batch = batch[:0]
				consumer.Commit()
			}
			k.logger.Info("received", slog.Any("dto", profileDto))
		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			if len(batch) > 0 {
				k.processBatch(ctx, batch)
				batch = batch[:0]
				consumer.Commit()
			}
			continue
		} else {
			k.logger.Error("error consuming message", slog.Any("error", err))
			// todo: нужно ли здесь делать коммит? И запроцессить текущий батч?
			// todo: поместить сообщение в dead letter queue
		}
	}

	return nil
}

func (k kafkaConsumer) processBatch(ctx context.Context, profilesDtos []ProfileDto) {
	k.logger.Info("start batch processing...")
	for _, dto := range profilesDtos {
		profile := model.Profile{
			ID:          dto.ID,
			Age:         dto.Age,
			Name:        dto.Name,
			Description: dto.Description,
			Photos:      dto.Photos,
		}

		err := k.storage.UpdateProfile(ctx, dto.Gender, profile)
		if err != nil {
			k.logger.Error("error updating profile in storage", slog.Any("error", err))
			// todo: publish to dead letter queue...
		}
	}
	k.logger.Info("finish batch processing")
}
