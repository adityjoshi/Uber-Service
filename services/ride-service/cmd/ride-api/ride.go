package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/db"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/handler"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/kafka"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/repository"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := db.Init(ctx); err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.RidePool.Close()

	producer := kafka.NewProducer()
	defer producer.Close()

	consumer := kafka.NewConsumer()
	defer consumer.Close()

	repo := repository.NewRideRepository(db.RidePool)
	svc := service.NewRideService(repo, producer)

	go consumer.Start(ctx, func(ctx context.Context, event kafka.RideMatchedEvent) error {
		return svc.UpdateRideWithDriver(ctx, event.RideId, event.DriverID)
	})
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	r.GET("/ready", readyHandler)

	rideHandler := handler.NewRideHandler(svc)
	rideHandler.RegisterRoutes(r)

	addr := getenv("HTTP_ADDR", ":8083")
	log.Printf("ride-service started listening on %s \n", addr)
	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatal(err)
		}

	}()

	<-ctx.Done()
	log.Println("ride service shutting down")

}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func readyHandler(c *gin.Context) {
	brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokers == "" {
		brokers = "kafka:9092"
	}
	first := strings.TrimSpace(strings.Split(brokers, ",")[0])
	if first == "" {
		c.JSON(503, gin.H{"error": "no kafka broker configured"})
		return
	}
	conn, err := net.DialTimeout("tcp", first, 2*time.Second)
	if err != nil {
		c.JSON(503, gin.H{"error": "no kafka broker configured"})
		return
	}
	_ = conn.Close()
	c.JSON(200, gin.H{"status": "kafka ok"})
}
