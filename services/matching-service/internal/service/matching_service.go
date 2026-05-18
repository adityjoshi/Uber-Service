package service

import (
	"context"
	"log"

	"github.com/adityjoshi/Uber-Service/services/matching-service/internal/dto"
	"github.com/adityjoshi/Uber-Service/services/matching-service/internal/kafka"
)

const (
	defaultSearchRadiuskm = 5.0
	distanceWeight        = 0.7
	ratingWeight          = 0.3
)

type MatchingService struct {
	locationClient *client.LocationClient
	producer       *kafka.Producer
}

func NewMatchingService(locationClient *client.LocationClient, producer *kafka.Producer) *MatchingService {
	return &MatchingService{
		locationClient: locationClient,
		producer:       producer,
	}
}

func (s *MatchingService) MathDriverForRide(ctx context.Context, event kafka.RideRequestedEvent) error {
	log.Printf("matching: finding rider for rideID=%s", event.RideID)

	drivers, err := s.locationClient.GetNearbyDrivers(ctx, event.PickupLatitude, event.PickupLongitude, defaultSearchRadiuskm)
	if err != nil {
		return err
	}
	if len(drivers) == 0 {
		log.Printf("matching: no suitable drivers for rideID=%s", event.RideID)
		return nil
	}

	best, found := findBestDriver(drivers)
	if !found {
		log.Printf("matching: no suitable driver for rideID=%s", event.RideID)
		return nil
	}

	matchedEvent := kafka.RideMatchedEvent{
		RideId:             event.RideID,
		RiderID:            event.RiderID,
		DriverID:           best.DriverID,
		DriverLatitude:     best.Latitude,
		DriverLongitude:    best.Longitude,
		DistanceToPickupKm: best.DistanceInKm,
	}
	if err := s.producer.PublishRideMatcher(ctx, matchedEvent); err != nil {
		return err
	}

	log.Printf("matching: ride.matched published - rideID=%s driverID=%s distance=%.2fkm", event.RideID, best.DriverID, best.DistanceInKm)
	return nil
}

func findBestDriver(drivers []dto.NearByDriverResponse) (dto.NearByDriverResponse, bool) {
	if len(drivers) == 0 {
		return dto.NearByDriverResponse{}, false
	}

	best := drivers[0]
	bestScore := score(best)

	for _, d := range drivers[1:] {
		if s := score(d); s > bestScore {
			bestScore = s
			best = d
		}
	}
	return best, true
}
