package kafka

import (
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

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
