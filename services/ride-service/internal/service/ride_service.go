package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/dto"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/kafka"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/model"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/repository"
	"github.com/go-playground/validator/v10/translations/id"
	"github.com/google/uuid"
)

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
