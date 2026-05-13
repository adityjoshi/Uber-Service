package kafka

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

type RideMatchedEvent struct {
	RideId             string  `json:"rideId"`
	RiderID            string  `json:"riderId"`
	DriverID           string  `json:"driverID"`
	DriverLatitude     float64 `json:"driverLatitude"`
	DriverLongitude    float64 `json:"driverLongitude"`
	DistanceToPickupKm float64 `json:"distanceToPickupKm"`
}
