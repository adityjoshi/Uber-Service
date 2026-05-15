package handler

import (
	"log"
	"net/http"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/dto"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/handler"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/service"
	"github.com/gin-gonic/gin"
)

type RideHandler struct {
	svc *service.RideService
}

func NewRideHandler(svc *service.RideService) *RideHandler {
	return &RideHandler{svc: svc}
}

func (h *RideHandler) RegisterRoutes(r gin.Engine) {

}

func (h *RideHandler) requestRide(c *gin.Context) {
	var req dto.RideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RiderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rider id is required"})
		return
	}
	log.Printf("handler: ride request received from the ride: %s", req.RiderID)

	resp, err := h.svc.RequestRide(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *RideHandler) getRide(c *gin.Context) {
	riderID := c.Param("rideID")
	rides, err := h.svc.GetRide(c.Request.Context(), riderID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, rides)
}

func (h *RideHandler) listByRider(c *gin.Context) {
	riderID := c.Param("riderID")

	rides, err := h.svc.ListByRider(c.Request.Context(), riderID)
	if err != nil {
		h.handleServiceError(c, err)
		return

	}
	c.JSON(http.StatusOK, rides)
}

func (h *RideHandler) startRide(c *gin.Context) {
	riderID := c.Param("riderID")

	resp, err := h.svc.StartRide(c.Request.Context(), riderID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)

}
