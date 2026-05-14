package dto

import (
	"time"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/model"
)

type RideRequest struct {
	RiderID       string  `json:"riderId"`
	PickupLat     float64 `json:"pickupLat"`
	PickupLong    float64 `json:"pickupLong"`
	PickupAddress string  `json:"pickupAddress"`
	DropLat       float64 `json:"dropLat"`
	DropLong      float64 `json:"dropLong"`
	DropAddress   string  `json:"dropAddress"`
}

type RideResponse struct {
	Id            string           `json:"id"`
	RiderId       string           `json:"riderId"`
	DriverId      string           `json:"driverId"`
	PickupLat     float64          `json:"pickupLat"`
	PickupLong    float64          `json:"pickupLong"`
	PickupAddress string           `json:"pickupAddress"`
	DropLat       float64          `json:"dropLat"`
	DropLong      float64          `json:"dropLong"`
	DropAddress   string           `json:"dropAddress"`
	Status        model.RideStatus `json:"status"`
	EstimatedFare float64          `json:"estimatedFare"`
	ActualFare    float64          `json:"actualFare"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
	StartedAt     *time.Time       `json:"startedAt"`
	CompletedAt   *time.Time       `json:"completedAt"`
}
