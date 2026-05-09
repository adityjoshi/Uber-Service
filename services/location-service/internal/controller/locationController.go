package controller

import "github.com/adityjoshi/Uber-Service/services/location-service/internal/services"

type LocationHandler struct {
	svc *services.LocationService
}

func NewLocationHandler(svc *services.LocationService) *LocationHandler {
	return &LocationHandler{svc: svc}
}
