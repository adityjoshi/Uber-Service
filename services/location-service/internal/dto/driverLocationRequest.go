package dto

// gives the driver's current location
type DriverLocationReq struct {
	DriverID string  `json:"driverId"`
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
}
