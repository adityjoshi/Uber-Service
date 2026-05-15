package handler

import (
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
