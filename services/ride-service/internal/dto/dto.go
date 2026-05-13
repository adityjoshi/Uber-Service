package dto

import "time"

type RideStatus string

const (
	RideStatusRequested      RideStatus = "REQUESTED"
	RideStatusMatching       RideStatus = "MATCHING"
	RideStatusAccepted       RideStatus = "ACCEPTED"
	RideStatusDriverArriving RideStatus = "DRIVER_ARIVING"
	RideStatusRideStarted    RideStatus = "RIDE_STARTED"
	RideStatusCompleted      RideStatus = "COMPLETED"
	RideStatusCancelled      RideStatus = "CANCELLED"
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
	Id            string     `json:"id"`
	RiderId       float64    `json:"riderId"`
	DriverId      float64    `json:"driverId"`
	PickupLat     float64    `json:"pickupLat"`
	PickupLong    float64    `json:"pickupLong"`
	PickupAddress string     `json:"pickupAddress"`
	DropLat       float64    `json:"dropLat"`
	DropLong      float64    `json:"dropLong"`
	DropAddress   string     `json:"dropAddress"`
	Status        RideStatus `json:"status"`
	EstimatedFare float64    `json:"estimatedFare"`
	ActualFare    float64    `json:"actualFare"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	StartedAt     *time.Time `json:"startedAt"`
	CompletedAt   *time.Time `json:"completedAt"`
}
