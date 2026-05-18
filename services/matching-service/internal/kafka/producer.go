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
	log.Printf("kafka prodcuer connected to %v", brokers)
	return &Producer{writer: writer}
}

func (p *Producer) PublishRideMatcher(ctx context.Context, event RideMatchedEvent) error {
	return p.publish(ctx, getenv("TOPIC_RIDE_MATCHED", "ride.matched"), event.RideId, event)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) publish(ctx context.Context, topic, key string, payload any) error {
	value, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("kafka Producer: matching service: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}); err != nil {
		return fmt.Errorf("kafka prodcuer: write to %s: %w", topic, err)
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
