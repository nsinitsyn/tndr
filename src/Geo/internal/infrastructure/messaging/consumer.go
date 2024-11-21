package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"
	"tinder-geo/internal/config"
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/infrastructure/messaging/dto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const BATCH_SIZE int = 10
const READ_MESSAGE_TIMEOUT_SEC = 5

type Service interface {
	UpdateProfile(ctx context.Context, gender model.Gender, profile model.Profile) error
}

type kafkaConsumer struct {
	config   config.MessagingConfig
	consumer *kafka.Consumer
	logger   *slog.Logger
	service  Service
}

func NewConsumer(config config.MessagingConfig, logger *slog.Logger, service Service) (kafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  config.Servers,
		"group.id":           config.Group,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return kafkaConsumer{}, err
	}

	return kafkaConsumer{config: config, consumer: consumer, logger: logger, service: service}, nil
}

func (k kafkaConsumer) MustSubscribe() {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": k.config.Servers})
	if err != nil {
		panic(err)
	}

	// Create topic if not exists
	_, err = admin.GetMetadata(&k.config.Topic, false, 5000)
	if err != nil {
		panic(err)
	}

	_, err = admin.CreateTopics(context.Background(), []kafka.TopicSpecification{{
		Topic:             k.config.Topic,
		NumPartitions:     1,
		ReplicationFactor: 1}})
	if err != nil {
		panic(err)
	}

	err = k.consumer.SubscribeTopics([]string{k.config.Topic}, nil)
	if err != nil {
		panic(err)
	}
}

func (k kafkaConsumer) StartConsume(ctx context.Context) error {
	k.logger.Info("consuming started")

	batch := make([]dto.ProfileDto, 0, BATCH_SIZE)

	for !errors.Is(ctx.Err(), context.Canceled) {
		msg, err := k.consumer.ReadMessage(READ_MESSAGE_TIMEOUT_SEC * time.Second)
		if err == nil {
			var profileDto dto.ProfileDto
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
				k.consumer.Commit()
			}
			k.logger.Info("received", slog.Any("dto", profileDto))
		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			if len(batch) > 0 {
				k.processBatch(ctx, batch)
				batch = batch[:0]
				k.consumer.Commit()
			}
			continue
		} else {
			k.logger.Error("error consuming message", slog.Any("error", err))
		}
	}

	return k.consumer.Close()
}

func (k kafkaConsumer) processBatch(ctx context.Context, profilesDtos []dto.ProfileDto) {
	k.logger.Info("start batch processing...")
	for _, dto := range profilesDtos {
		profile := model.Profile{
			ID:          dto.ID,
			Age:         dto.Age,
			Name:        dto.Name,
			Description: dto.Description,
			Photos:      dto.Photos,
		}

		err := k.service.UpdateProfile(ctx, dto.Gender, profile)
		if err != nil {
			k.logger.Error("error updating profile", slog.Any("error", err))
			// todo: publish to dead letter queue...
		}
	}
	k.logger.Info("finish batch processing")
}
