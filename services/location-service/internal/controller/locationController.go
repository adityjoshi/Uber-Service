package controller

import (
	"errors"
	"net/http"

	"github.com/adityjoshi/Uber-Service/services/location-service/internal/dto"
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

func (h *LocationHandler) UpdateDriverLocation(c *gin.Context) {
	var req dto.DriverLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	err := h.svc.UpdateDriverLocation(c.Request.Context(), req.DriverID, req.Lat, req.Long)
	if err != nil {
		if errors.Is(err, services.ErrDriverIdRequired) {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "driver location updated"})
}
