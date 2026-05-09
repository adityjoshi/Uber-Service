package controller

import (
	"github.com/adityjoshi/Uber-Service/services/location-service/internal/services"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	svc *services.LocationService
}

func NewLocationHandler(svc *services.LocationService) *LocationHandler {
	return &LocationHandler{svc: svc}
}

func (h *LocationHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/locations")

	{
		g.POST("/driver/update", h.UpdateDriverLocation)
		g.GET("/driver/nearby", h.GetNearbyDriver)
		g.DELETE("/driver/:driverID", h.RemoveDriver)
	}
}
