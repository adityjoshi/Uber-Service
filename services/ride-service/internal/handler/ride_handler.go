package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/dto"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/service"
	"github.com/gin-gonic/gin"
)

type RideHandler struct {
	svc *service.RideService
}

func NewRideHandler(svc *service.RideService) *RideHandler {
	return &RideHandler{svc: svc}
}

func (h *RideHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/rides")
	{
		v1.POST("/request", h.requestRide)
		v1.GET("/:rideId", h.getRide)
		v1.GET("/rider/:riderId", h.listByRider)
		v1.PUT("/:rideId/start", h.startRide)
		v1.PUT("/:rideId/complete", h.completeRide)
		v1.PUT("/:rideId/cancel", h.cancelRide)
	}
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
	riderID := c.Param("rideId")
	rides, err := h.svc.GetRide(c.Request.Context(), riderID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, rides)
}

func (h *RideHandler) listByRider(c *gin.Context) {
	riderID := c.Param("riderId")

	rides, err := h.svc.ListByRider(c.Request.Context(), riderID)
	if err != nil {
		h.handleServiceError(c, err)
		return

	}
	c.JSON(http.StatusOK, rides)
}

func (h *RideHandler) startRide(c *gin.Context) {
	riderID := c.Param("rideId")

	resp, err := h.svc.StartRide(c.Request.Context(), riderID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)

}

func (h *RideHandler) completeRide(c *gin.Context) {
	rideID := c.Param("rideId")

	resp, err := h.svc.CompleteRide(c.Request.Context(), rideID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *RideHandler) cancelRide(c *gin.Context) {
	rideID := c.Param("rideId")

	resp, err := h.svc.CancelRide(c.Request.Context(), rideID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *RideHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNoRideFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrInvalidStatus):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		log.Printf("handler internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

}
