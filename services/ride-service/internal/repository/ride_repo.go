package repository

import (
	"context"
	"fmt"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RideRepository struct {
	db *pgxpool.Pool
}

func NewRideRepository(db *pgxpool.Pool) *RideRepository {
	return &RideRepository{db: db}
}

func (r *RideRepository) Save(ctx context.Context, ride *model.Ride) error {
	query := `
	INSERT INTO rides (
		id, rider_id, driver_id, pickup_latitude, pickup_longitude, pickup_address,
		drop_latitude, drop_longitude, drop_address, status,
		estimated_fare, actual_fare, created_at, updated_at, started_at, completed_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
	)
	ON CONFLICT (id) DO UPDATE SET

	driver_id = EXCLUDED.driver_id,
	status = EXCLUDED.status,
	estimated_fare = EXCLUDED.estimated_fare,
	actual_fare = EXCLUDED.actual_fare,
	started_at = EXCLUDED.started_at,
	completed_at = EXCLUDED.completed_at
	`
	_, err := r.db.Exec(ctx, query, ride.ID, ride.RiderID, ride.DriverID, ride.PickupLatitude, ride.PickupLongitude, ride.PickupAddress,
		ride.DropLatitude, ride.DropLongitude, ride.DropAddress, ride.Status, ride.EstimatedFare, ride.ActualFare, ride.CreatedAt, ride.UpdatedAt, ride.StartedAt, ride.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("repository: save ride: %w", err)
	}

	return nil

}

func (r *RideRepository) FindById(ctx context.Context, id string) (*model.Ride, error) {

	query := `
	SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, pickup_address, drop_latitude, drop_longitude, drop_address, status, estimated_fare,actual_fare,created_at,updated_at,started_at,completed_at
	FROM rides
	WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	ride, err := scanRide(row)
	if err != nil {
		return nil, fmt.Errorf("repository: find ride by id: %w", err)
	}
	return ride, nil
}

func (r *RideRepository) FindByRiderId(ctx context.Context, riderId string) ([]*model.Ride, error) {
	query := `
		SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, pickup_address, drop_latitude, drop_longitude, drop_address, status, estimated_fare,actual_fare,created_at,updated_at,started_at,completed_at
	FROM rides where rider_id = $1 ORDER BY created_at DESC

	`
	rows, err := r.db.Query(ctx, query, riderId)
	if err != nil {
		return nil, fmt.Errorf("repository: find rides by rider: %w", err)
	}
	defer rows.Close()

	var rides []*model.Ride
	for rows.Next() {
		ride, err := scanRide(rows)
		if err != nil {
			return nil, fmt.Errorf("repository: scan ride: %w", err)
		}
		rides = append(rides, ride)
	}
	return rides, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRide(s scanner) (*model.Ride, error) {
	var ride model.Ride
	err := s.Scan(&ride.ID, &ride.RiderID, &ride, ride.DriverID, &ride.PickupLatitude, &ride.PickupLongitude, &ride.PickupAddress, &ride.DropLatitude, &ride.DropLongitude,
		&ride.DropAddress, &ride.Status, &ride.EstimatedFare, &ride.ActualFare, &ride.CreatedAt, &ride.UpdatedAt, &ride.StartedAt, &ride.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ride, nil
}
