package services

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/adityjoshi/Uber-Service/services/location-service/internal/dto"
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

func (s *LocationService) FindNearbyDriver(
	ctx context.Context,
	lat float64,
	long float64,
	radiusInKm float64,
) ([]dto.NearByDriverResponse, error) {

	results, err := s.rdb.GeoSearchLocation(
		ctx,
		driversGeoKey,
		&redis.GeoSearchLocationQuery{
			GeoSearchQuery: redis.GeoSearchQuery{
				Longitude:  long,
				Latitude:   lat,
				Radius:     radiusInKm,
				RadiusUnit: "km",
				Sort:       "ASC",
				Count:      10,
			},
			WithCoord: true,
			WithDist:  true,
		},
	).Result()

	if err != nil {
		return nil, err
	}

	drivers := make([]dto.NearByDriverResponse, 0, len(results))

	for _, result := range results {
		driver := dto.NearByDriverResponse{
			DriverID: result.Name,
			Distance: result.Dist,
			Lat:      result.Latitude,
			Long:     result.Longitude,
		}
		drivers = append(drivers, driver)
	}
	return drivers, nil
}

func (s *LocationService) RemoveDriver(ctx context.Context, driverID string) error {
	if driverID == "" {
		return errors.New("driver id is required")
	}
	return s.rdb.ZRem(ctx, driversGeoKey, driverID).Err()
}

func (s *LocationService) Validate() error {
	if s == nil || s.rdb == nil {
		return fmt.Errorf("location is not init")
	}
	return nil
}
