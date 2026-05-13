package repository

import "github.com/jackc/pgx/v5/pgxpool"

type RideRepository struct {
	db *pgxpool.Pool
}

func NewRideRepository(db *pgxpool.Pool) *RideRepository {
	return &RideRepository{db: db}
}
