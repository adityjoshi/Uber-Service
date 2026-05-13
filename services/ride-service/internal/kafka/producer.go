package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	brokers := strings.Split(getenv("KAFKA_BROKERS", "kafka:9092"), ",")

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	log.Printf("kakfa Producer connected to %v", brokers)
	return &Producer{writer: writer}
}

func (p *Producer) PublishRideRequested(ctx context.Context, event RideRequestedEvent) error {
	return p.publish(ctx, getenv("TOPIC_RIDE_REQUESTED", "ride.requested"), event.RideID, event)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) publish(ctx context.Context, topic, key string, payload any) error {
	value, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("kafka producer: Marshal: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})

	if err != nil {
		return fmt.Errorf("kafka producer: writer to %s: %w", topic, err)
	}
	log.Printf("kafka producer: published to %s key=%s", topic, key)
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
