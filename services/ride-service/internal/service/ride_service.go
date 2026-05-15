package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/dto"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/kafka"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/model"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrInvalidStatus = errors.New("ride not found")
var ErrNoRideFound = errors.New("no rides found for the rider")

type RideService struct {
	repo     *repository.RideRepository
	producer *kafka.Producer
}

func NewRideService(repo *repository.RideRepository, producer *kafka.Producer) *RideService {
	return &RideService{repo: repo, producer: producer}
}

func (s *RideService) RequestRide(ctx context.Context, req dto.RideRequest) (*dto.RideResponse, error) {
	log.Printf("A new ride request coming from rider: %s", req.RiderID)
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location", err)
		return nil, err
	}

	now := time.Now().In(loc)

	ride := &model.Ride{
		ID:              uuid.NewString(),
		RiderID:         req.RiderID,
		PickupLatitude:  req.PickupLat,
		PickupLongitude: req.PickupLong,
		PickupAddress:   req.PickupAddress,
		DropLatitude:    req.DropLat,
		DropLongitude:   req.DropLong,
		DropAddress:     req.DropAddress,
		Status:          model.RideStatusRequested,
		EstimatedFare:   calculateFare(req.PickupLatitude, req.PickupLong, req.DropLat, req.DropLong),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.repo.Save(ctx, ride); err != nil {
		return nil, fmt.Errorf("service: request ride: %w", err)
	}
	event := kafka.RideRequestedEvent{
		RideID:          ride.ID,
		RiderID:         ride.RiderID,
		PickupLatitude:  ride.PickupLatitude,
		PickupLongitude: ride.PickupLongitude,
		PickupAddress:   ride.PickupAddress,
		DropLatitude:    ride.DropLatitude,
		DropLongitude:   ride.DropLongitude,
		DropAddress:     ride.DropAddress,
	}

	if err := s.producer.PublishRideRequested(ctx, event); err != nil {
		log.Printf("service: failed to publish the ride.requested for rideId=%s: %v", err)
	}
	log.Printf("service: ride.requested published for rideID=%s", ride.ID)

	ride.Status = model.RideStatusMatching
	if err := s.repo.Save(ctx, ride); err != nil {
		return nil, fmt.Errorf("service: update ride to matching: %w", err)
	}
	return mapToResponse(ride), err
}

/*
* It is called by kafka consumer when the ride.matched is received
* */

func (s *RideService) UpdateRideWithDriver(ctx context.Context, rideID, driverID string) error {
	ride, err := s.findOrNotFound(ctx, rideID)
	if err != nil {
		return err
	}
	ride.DriverId = &driverID
	ride.Status = model.RideStatusAccepted

	if err := s.repo.Save(ctx, ride); err != nil {
		return fmt.Errorf("service: update ride with the driver: %w", err)
	}
	return nil
}

/*
* changes the accepted -> ride started
* */

func (s *RideService) StartRide(ctx context.Context, rideID string) (*dto.RideResponse, error) {
	ride, err := s.findOrNotFound(ctx, rideID)
	if err != nil {
		return nil, err
	}

	if ride.Status != model.RideStatusAccepted {
		return nil, fmt.Errorf("service: start ride error: %w", err)
	}
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location", err)
		return nil, err
	}

	now := time.Now().In(loc)
	ride.Status = model.RideStatusRideStarted
	ride.StartedAt = &now

	if err := s.repo.Save(ctx, ride); err != nil {
		return nil, fmt.Errorf("service: ride start error: %w", err)
	}
	return mapToResponse(ride), nil
}

/*
* Complete ride changes the ride started -> completed
* */

func (s *RideService) CompleteRide(ctx context.Context, rideID string) (*dto.RideResponse, error) {
	ride, err := s.findOrNotFound(ctx, rideID)
	if err != nil {
		return nil, err
	}
	if ride.Status != model.RideStatusRideStarted {
		return nil, fmt.Errorf("%w: cannot complete, current status is %s", ErrInvalidStatus, ride.Status)
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)

	ride.Status = model.RideStatusCompleted
	ride.CompletedAt = &now
	ride.ActualFare = ride.EstimatedFare

	if err != s.repo.Save(ctx, ride); err != nil {
		return nil, fmt.Errorf("service: complete ride error: %w", err)
	}
	return mapToResponse(ride), nil
}

func (s *RideService) CancelRide(ctx context.Context, rideID string) (*dto.RideResponse, error) {
	ride, err := s.findOrNotFound(ctx, rideID)
	if err != nil {
		return nil, err
	}

	ride.Status = model.RideStatusCancelled

	if err := s.repo.Save(ctx, ride); err != nil {
		return nil, fmt.Errorf("service: cancel ride: %w", err)
	}
	return mapToResponse(ride), nil
}

func (s *RideService) GetRide(ctx context.Context, rideID string) (*dto.RideResponse, error) {
	ride, err := s.findOrNotFound(ctx, rideID)
	if err != nil {
		return nil, err
	}
	return mapToResponse(ride), nil
}

/*
*List all the rides for the rider. Newest rides will be seen first
* */

func (s *RideService) ListByRider(ctx context.Context, riderID string) ([]*dto.RideResponse, error) {
	rides, err := s.findOrNotFound(ctx, riderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrNoRideFound, riderID)
		}
		return nil, fmt.Errorf("service: fetch rides error: %w", err)
	}
	responses := make([]*dto.RideResponse, len(rides))
	for i, r := range rides {
		responses[i] = mapToResponse(r)
	}
	return responses, nil

}

func (s *RideService) findOrNotFound(ctx context.Context, rideID string) (*model.Ride, error) {
	ride, err := s.repo.FindById(ctx, rideID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrNoRideFound, rideID)
		}
		return nil, fmt.Errorf("service: fetch ride: %w", err)
	}
	return ride, nil
}

// Base fare: ₹50 + ₹12/km, rounded to 2 decimal places.
func calculateFare(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Pow(math.Sin(dLon/2), 2)

	distanceKm := earthRadiusKm * 2 * math.Asin(math.Sqrt(a))

	fare := 50 + (distanceKm * 12)
	return math.Round(fare*100) / 100
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func mapToResponse(r *model.Ride) *dto.RideResponse {
	resp := &dto.RideResponse{
		Id:            r.ID,
		RiderId:       r.RiderID,
		PickupLat:     r.PickupLatitude,
		PickupLong:    r.PickupLongitude,
		PickupAddress: r.PickupAddress,
		DropLat:       r.DropLatitude,
		DropLong:      r.DropLongitude,
		DropAddress:   r.DropAddress,
		Status:        r.Status,
		EstimatedFare: r.EstimatedFare,
		ActualFare:    r.ActualFare,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		StartedAt:     r.StartedAt,
		CompletedAt:   r.CompletedAt,
	}
	if r.DriverID != nil {
		resp.DriverId = *r.DriverID
	}
	return resp
}
