package services

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type LocationService struct {
	rdb *redis.Client
}

func NewLocationService(rdb *redis.Client) *LocationService {
	return &LocationService{rdb: rdb}
}

const driversGeoKey = "drivers:location"

func (s *LocationService) UpdateDriverLocation(
	ctx context.Context,
	driverID string,
	lat float64,
	long float64) error {
	if driverID == "" {
		return errors.New("driverid is required")
	}

	return s.rdb.GeoAdd(ctx, driversGeoKey, &redis.GeoLocation{
		Name:      driverID,
		Longitude: long,
		Latitude:  lat,
	}).Err()
}
