package services

import "github.com/redis/go-redis/v9"

type LocationService struct {
	rdb *redis.Client
}

func NewLocationService(rdb *redis.Client) *LocationService {
	return &LocationService{rdb: rdb}
}
