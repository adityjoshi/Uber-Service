package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {

	ctx := context.Background()

	if err := db.Init(ctx); err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.RidePool.Close()

	r := gin.Default()

	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	r.GET("/ready", readyHandler)

	addr := getenv("HTTP_ADDR", ":8083")
	log.Println("ride-service started listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}

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
