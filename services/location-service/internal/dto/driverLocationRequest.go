package dto

// gives the driver's current location
type DriverLocationReq struct {
	DriverID string
	Lat      float64
	Long     float64
}
