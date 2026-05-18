package kafka

/*
*
* RideRequestedEvent is published to topic ride.requested
* after a new ride is saved. The matching service will consume this.
*
* */

type RideRequestedEvent struct {
	RideID          string  `json:"rideId"`
	RiderID         string  `json:"riderId"`
	PickupLatitude  float64 `json:"pickupLatitude"`
	PickupLongitude float64 `json:"pickupLongitude"`
	PickupAddress   string  `json:"pickupAddress"`
	DropLatitude    float64 `json:"dropLatitude"`
	DropLongitude   float64 `json:"dropLongitude"`
	DropAddress     string  `json:"dropAddress"`
}

/*
*
* RideMatchedEvent is consumed from the topic ride.matched
*
* published by the matching service when a driver is assigned
* */

type RideMatchedEvent struct {
	RideId             string  `json:"rideId"`
	RiderID            string  `json:"riderId"`
	DriverID           string  `json:"driverID"`
	DriverLatitude     float64 `json:"driverLatitude"`
	DriverLongitude    float64 `json:"driverLongitude"`
	DistanceToPickupKm float64 `json:"distanceToPickupKm"`
}
