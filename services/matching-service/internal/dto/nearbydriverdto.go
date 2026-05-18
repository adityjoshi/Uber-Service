package dto

type NearByDriverResponse struct {
	DriverID string  `json:"driverId"`
	Lat      float64 `json:"latitude"`
	Long     float64 `json:"longitude"`
	Distance float64 `json:"distanceInKm"`
}
