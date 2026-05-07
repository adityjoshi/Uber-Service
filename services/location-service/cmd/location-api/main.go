package main

import (
	"context"
	"log"
	"time"

	"github.com/adityjoshi/Uber-Service/services/location-service/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCfg := config.LoadRedisFromEnv()
	rdb := config.NewRedisClient(redisCfg)
	defer rdb.Close()

	if err := config.PingRedis(ctx, rdb); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	log.Println("Redis Connection Success")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	log.Fatal(router.Run(":8082"))
}
