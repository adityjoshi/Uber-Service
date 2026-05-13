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
