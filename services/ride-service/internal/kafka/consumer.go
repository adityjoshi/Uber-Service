package kafka

import (
	"context"
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
