package service

import (
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/kafka"
	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/repository"
)

type RideService struct {
	repo     *repository.RideRepository
	producer *kafka.Producer
}

func NewRideService(repo *repository.RideRepository, producer *kafka.Producer) *RideService {
	return &RideService{repo: repo, producer: producer}
}
