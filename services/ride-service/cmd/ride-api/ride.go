package main

import (
	"context"
	"log"

	"github.com/adityjoshi/Uber-Service/services/ride-service/internal/db"
)

func main() {

	ctx := context.Background()

	if err := db.Init(ctx); err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
}
