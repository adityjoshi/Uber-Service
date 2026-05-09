package controller

import (
	"errors"
	"net/http"
	"strconv"

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
		if errors.Is(err, services.ErrDriverIDRequired) {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "driver location updated"})
}

func (h *LocationHandler) GetNearbyDriver(c *gin.Context) {
	latStr := c.Query("latitude")
	longStr := c.Query("longitude")

	if latStr == "" || longStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and long is required"})
		return
	}

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
		return
	}
	longitude, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
		return
	}

	radius := 5.0

	if rs := c.Query("radius"); rs != "" {
		radius, err := strconv.ParseFloat(rs, 64)
		if err != nil || radius <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid radius"})
			return
		}
	}

	drivers, err := h.svc.FindNearbyDriver(c.Request.Context(), latitude, longitude, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, drivers)
}
