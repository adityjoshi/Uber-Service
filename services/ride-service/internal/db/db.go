package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var RidePool *pgxpool.Pool

func Init(ctx context.Context) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("POSTGRES_USER", "rideuser"),
		getenv("POSTGRES_PASSWORD", "ridepassword"),
		getenv("POSTGRES_DB", "ridedb"),
		getenv("DB_SSLMODE", "disable"),
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("db : connect: %w ", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("db: ping %w", err)
	}
	RidePool = pool

	log.Println("connected to the postgres")
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
