package service

import (
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/kafka"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/repository"
)

type RideService struct {
	repo     *repository.RideRepository
	producer *kafka.Producer
}
