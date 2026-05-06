package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host string
	Port int
}

func LoadRedisFromEnv() RedisConfig {
	portStr := os.Getenv("REDIS_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 6379
	}
	return RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: port,
	}
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func NewRedisClient(cfg RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})
}

func PingRedis(ctx context.Context, rdb *redis.Client) error {
	return rdb.Ping(ctx).Err()
}

func getEnvIntDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
