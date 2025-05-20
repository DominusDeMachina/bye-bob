package repository

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gfurduy/byebob/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Global database pool instance
var (
	globalDBPool *DBPool
	poolOnce     sync.Once
	poolMutex    sync.RWMutex
)

// DBPool is a wrapper around pgxpool.Pool with common functionality
type DBPool struct {
	Pool *pgxpool.Pool
	Cfg  *config.Config
}

// PoolConfig defines configuration for the database connection pool
type PoolConfig struct {
	MaxConns         int32         // Maximum number of connections in the pool (default: 10)
	MinConns         int32         // Minimum number of connections in the pool (default: 2)
	MaxConnLifetime  time.Duration // Maximum lifetime of connections (default: 1 hour)
	MaxConnIdleTime  time.Duration // Maximum idle time of connections (default: 30 minutes)
	HealthCheckPeriod time.Duration // Period between health checks (default: 1 minute)
	ConnectTimeout   time.Duration // Timeout for establishing connections (default: 5 seconds)
	MaxRetries       int           // Maximum number of connection retries (default: 5)
	RetryDelay       time.Duration // Initial delay between retries (default: 3 seconds)
}

// DefaultPoolConfig returns the default connection pool configuration
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MaxConns:         10,
		MinConns:         2,
		MaxConnLifetime:  time.Hour,
		MaxConnIdleTime:  30 * time.Minute,
		HealthCheckPeriod: time.Minute,
		ConnectTimeout:   5 * time.Second,
		MaxRetries:       5,
		RetryDelay:       3 * time.Second,
	}
}

// NewDBPool creates a new database connection pool with default settings
func NewDBPool(cfg *config.Config) (*DBPool, error) {
	return NewDBPoolWithConfig(cfg, DefaultPoolConfig())
}

// NewDBPoolWithConfig creates a new database connection pool with custom configuration
func NewDBPoolWithConfig(cfg *config.Config, poolCfg *PoolConfig) (*DBPool, error) {
	var pool *pgxpool.Pool
	var connectConfig *pgxpool.Config
	var err error

	// Parse the connection string
	connectConfig, err = pgxpool.ParseConfig(cfg.PostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Apply pool configuration
	connectConfig.MaxConns = poolCfg.MaxConns
	connectConfig.MinConns = poolCfg.MinConns
	connectConfig.MaxConnLifetime = poolCfg.MaxConnLifetime
	connectConfig.MaxConnIdleTime = poolCfg.MaxConnIdleTime
	connectConfig.HealthCheckPeriod = poolCfg.HealthCheckPeriod

	// Attempt to connect with retries
	retryDelay := poolCfg.RetryDelay
	var lastError error

	for i := 0; i < poolCfg.MaxRetries; i++ {
		// Create context with timeout for connection attempt
		ctx, cancel := context.WithTimeout(context.Background(), poolCfg.ConnectTimeout)
		
		// Create the connection pool
		pool, err = pgxpool.NewWithConfig(ctx, connectConfig)
		
		// Test the connection with a ping
		if err == nil {
			err = pool.Ping(ctx)
		}
		
		cancel() // Cancel the context regardless of the outcome
		
		if err == nil {
			// Connection successful
			log.Printf("Successfully connected to PostgreSQL database (attempt %d/%d)", i+1, poolCfg.MaxRetries)
			break
		}
		
		// Connection failed
		lastError = err
		
		// Close the pool if it was created
		if pool != nil {
			pool.Close()
			pool = nil
		}
		
		// If this is not the last attempt, wait and retry
		if i < poolCfg.MaxRetries-1 {
			log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, poolCfg.MaxRetries, err)
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
			
			// Exponential backoff: increase delay for next retry
			retryDelay = retryDelay * 2
		}
	}

	// If all connection attempts failed
	if pool == nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", poolCfg.MaxRetries, lastError)
	}

	return &DBPool{
		Pool: pool,
		Cfg:  cfg,
	}, nil
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

// HealthCheck performs a health check on the database connection
func (db *DBPool) HealthCheck(ctx context.Context) error {
	// Create a timeout context if one wasn't provided
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Run a simple query to verify the connection
	var result int
	err := db.Pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	
	// Check if the result is as expected
	if result != 1 {
		return fmt.Errorf("database health check returned unexpected result: %d", result)
	}
	
	return nil
}

// GetPool returns the underlying pgxpool.Pool
func (db *DBPool) GetPool() *pgxpool.Pool {
	return db.Pool
}

// AcquireConn acquires a connection from the pool with context
func (db *DBPool) AcquireConn(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection from pool: %w", err)
	}
	return conn, nil
}

// WithAcquire executes the given function with an acquired connection and releases it afterward
func (db *DBPool) WithAcquire(ctx context.Context, fn func(*pgxpool.Conn) error) error {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection from pool: %w", err)
	}
	defer conn.Release()
	
	return fn(conn)
}

// Stats returns statistics about the connection pool
func (db *DBPool) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}

// ResetStats resets the statistics of the connection pool
func (db *DBPool) ResetStats() {
	// Note: pgxpool.Stat doesn't have a Reset method in v5, so we can't reset stats directly
	// This is a placeholder for future implementations if needed
}

// InitGlobalDBPool initializes the global database pool with the given configuration
// It ensures the pool is only initialized once
func InitGlobalDBPool(cfg *config.Config) (*DBPool, error) {
	var err error
	
	poolOnce.Do(func() {
		log.Println("Initializing global database connection pool...")
		globalDBPool, err = NewDBPool(cfg)
	})
	
	return globalDBPool, err
}

// GetGlobalDBPool returns the global database pool instance
// If it hasn't been initialized, it returns nil
func GetGlobalDBPool() *DBPool {
	poolMutex.RLock()
	defer poolMutex.RUnlock()
	return globalDBPool
}

// CloseGlobalDBPool closes the global database pool if it exists
func CloseGlobalDBPool() {
	poolMutex.Lock()
	defer poolMutex.Unlock()
	
	if globalDBPool != nil {
		globalDBPool.Close()
		globalDBPool = nil
		log.Println("Global database pool closed")
	}
} 