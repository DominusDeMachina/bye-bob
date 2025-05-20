package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gfurduy/byebob/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool is a wrapper around pgxpool.Pool with common functionality
type DBPool struct {
	Pool *pgxpool.Pool
}

// NewDBPool creates a new database connection pool
func NewDBPool(cfg *config.Config) (*DBPool, error) {
	// Create a connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.PostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Set pool configuration
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	// Create a connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return &DBPool{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DBPool) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("Database connection pool closed")
	}
}

// Ping checks if the database is reachable
func (db *DBPool) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// GetPool returns the underlying pgxpool.Pool
func (db *DBPool) GetPool() *pgxpool.Pool {
	return db.Pool
} 