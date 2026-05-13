package kafka

import (
	"log"
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
