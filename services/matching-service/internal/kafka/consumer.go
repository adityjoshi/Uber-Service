package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

type RideRequestHandler func(ctx context.Context, event RideRequestHandler) error

func NewConsumer() *Consumer {
	brokers := strings.Split(getenv("KAFKA_BROKERS", "kafka_9092"), ",")
	topic := getenv("TOPIC_RIDE_REQUESTED", "ride.requested")
	group := getenv("KAFKA_CONSUMER_GROUP", "matching-service-group")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        group,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	log.Printf("kafka consumer listening on the topic=%s and group=%s", topic, group)
	return &Consumer{reader: reader}
}

func (c *Consumer) Start(ctx context.Context, handler RideRequestHandler) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("kafka consumer: context cancelled, shutting down")
				return
			}
			log.Printf("kafka consumer: fetch error: %v", err)
			continue
		}
		var event RideRequestedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("kafka consumer: unmarshal error: %v — skipping message", err)
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		if err := handler(ctx, event); err != nil {
			log.Printf("kafka consumer: handler error for rideId=%s: %v — will retry", event.RideID, err)
			continue // do not commit — message will be redelivered
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("kafka consumer: commit error: %v", err)
		}

		log.Printf("kafka consumer: processed ride.requested rideId=%s", event.RideID)
	}
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("kafka consumer: close: %w", err)
	}
	return nil
}
