package database

import (
	"context"
	"log"
	"time"

	"github.com/gfurduy/byebob/config"
)

// Initialize sets up the database connection and returns a DB instance
func Initialize(cfg *config.Config) (*DB, error) {
	// Attempt to connect to the database with retries
	var db *DB
	var err error
	maxRetries := 5
	retryDelay := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = New(cfg)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay = retryDelay * 2
		}
	}

	if err != nil {
		return nil, err
	}

	// Verify the connection by running a simple query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int
	err = db.Pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Database initialization completed successfully")
	
	return db, nil
}

// HealthCheck performs a simple database query to check health
func (db *DB) HealthCheck(ctx context.Context) error {
	var result int
	return db.Pool.QueryRow(ctx, "SELECT 1").Scan(&result)
} 