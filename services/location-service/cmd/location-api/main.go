package main

import (
	"context"
	"log"
	"time"

	"github.com/adityjoshi/Uber-Service/services/location-service/internal/config"
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
}
