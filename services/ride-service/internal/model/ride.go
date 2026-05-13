package model

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

type Ride struct {
	ID              string     `db:"id"`
	RiderID         string     `db:"rider_id"`
	DriverID        *string    `db:"driver_id"`
	PickupLatitude  float64    `db:"pickup_latitude"`
	PickupLongitude float64    `db:"pickup_longitude"`
	PickupAddress   string     `db:"pickup_address"`
	DropLatitude    float64    `db:"drop_latitude"`
	DropLongitude   float64    `db:"drop_longitude"`
	DropAddress     string     `db:"drop_address"`
	Status          RideStatus `db:"status"`
	EstimatedFare   float64    `db:"estimated_fare"`
	ActualFare      float64    `db:"actual_fare"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	StartedAt       *time.Time `db:"started_at"`
	CompletedAt     *time.Time `db:"completed_at"`
}
