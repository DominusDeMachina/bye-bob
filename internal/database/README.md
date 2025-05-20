# Database Package

This package provides a PostgreSQL connection management system with connection pooling using pgxpool. It's designed to work with both local PostgreSQL instances and Supabase PostgreSQL.

## Components

- **db.go** - Core database functionality including connection pool setup
- **init.go** - Database initialization with retry logic and health checks

## Usage

```go
package main

import (
	"context"
	"log"
	
	"github.com/gfurduy/byebob/config"
	"github.com/gfurduy/byebob/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Initialize database with retry logic
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	// Use the connection pool
	var result int
	err = db.Pool.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	
	log.Println("Database connection successful!")
}
```

## Configuration

Database connection settings are managed through the `config` package and are loaded from environment variables. Key settings include:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=byebob
DB_SSLMODE=disable

# For Supabase
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_API_KEY=your-api-key
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key
```

When `SUPABASE_URL` and `SUPABASE_SERVICE_KEY` are provided, the connection will be made to Supabase instead of a local PostgreSQL instance.

## Connection Pool Settings

The connection pool is configured with the following default settings:

- Max connections: 10
- Min connections: 2
- Max connection lifetime: 1 hour
- Max connection idle time: 30 minutes
- Health check period: 1 minute

These settings can be adjusted in `db.go` based on the application's needs and expected load.

## Error Handling

The package includes retry logic for database connections and wraps errors with additional context to make debugging easier.

## Health Checks

The `HealthCheck` method can be used to verify database connectivity:

```go
err := db.HealthCheck(context.Background())
if err != nil {
    log.Printf("Database health check failed: %v", err)
}
```

This is useful for monitoring the health of the database connection in production environments.

## Repository Pattern

This database package is designed to be used with the repository pattern, as defined in the `internal/repository` package. The connection pool should be passed to repository implementations to perform database operations. 