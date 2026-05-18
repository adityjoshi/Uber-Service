package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityjoshi/Uber-Service/services/matching-service/internal/client"
	"github.com/adityjoshi/Uber-Service/services/matching-service/internal/kafka"
	"github.com/adityjoshi/Uber-Service/services/matching-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Kafka
	producer := kafka.NewProducer()
	defer producer.Close()

	consumer := kafka.NewConsumer()
	defer consumer.Close()

	// Service
	locationClient := client.NewLocationClient()
	svc := service.NewMatchingService(locationClient, producer)

	// Kafka consumer — runs in background until ctx is cancelled
	go consumer.Start(ctx, func(ctx context.Context, event kafka.RideRequestedEvent) error {
		return svc.MathDriverForRide(ctx, event)
	})

	// HTTP — health + ready probes only (no business endpoints)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	addr := getenv("HTTP_ADDR", ":8084")
	log.Printf("matching-service listening on %s", addr)

	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("matching-service shutting down")
	time.Sleep(500 * time.Millisecond) // let in-flight messages drain
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
