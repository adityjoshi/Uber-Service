package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

type RideMatchHandler func(ctx context.Context, event RideMatchedEvent) error

func NewConsumer() *Consumer {
	brokers := strings.Split(getenv("KAFKA_BROKERS", "kafka:9092"), ",")
	topic := getenv("TOPIC_RIDE_MATCHED", "ride.matched")
	group := getenv("KAFKA_CONSUMER_GROUP", "ride-service-group")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        group,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})

	log.Printf("kafka consumer listening on topic=%s group=%s", topic, group)
	return &Consumer{reader: reader}
}

func (c *Consumer) Start(ctx context.Context, handler RideMatchHandler) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("kafka consumer: context cancelled, shutting down")
				return
			}
			log.Printf("kafka consumer fetch error: %v", err)
			continue
		}
		var event RideMatchedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("kafka consumer: unmarshal error: %v", err)
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		if err := handler(ctx, event); err != nil {
			log.Printf("kafka cosnumer: handler error for rideID=%s: %v - will retry", event.RideId, err)
			continue
		}
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("kafka consumer: commit error: %v", err)
		}

		log.Printf("kafka consumer: processed ride.matched rideId=%s driverID=%s", event.RideId, event.DriverID)
	}
}
