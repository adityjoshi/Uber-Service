package service

import "github.com/adityjoshi/Uber-Service/services/matching-service/internal/kafka"

const (
	defaultSearchRadiuskm = 5.0
	distanceWeight        = 0.7
	ratingWeight          = 0.3
)

type MatchingService struct {
	locationClient *client.LocationClient
	producer       *kafka.Producer
}

func NewMatchingService(locationClient *client.LocationClient, producer *kafka.Producer) *MatchingService {
	return &MatchingService{
		locationClient: locationClient,
		producer:       producer,
	}
}
