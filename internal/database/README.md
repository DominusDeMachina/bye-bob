# Database Package

This package provides a PostgreSQL connection management system with connection pooling using pgxpool. It's designed to work with both local PostgreSQL instances and Railway.com PostgreSQL.

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
# Individual connection parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=byebob
DB_SSLMODE=disable

# Railway.com (recommended)
RAILWAY_DB_URL=postgresql://postgres:password@containers-us-west-xxx.railway.app:5432/railway
```

When `RAILWAY_DB_URL` is provided, the connection will use the Railway.com PostgreSQL instance. Otherwise, it will use the individual connection parameters.

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

# Database Connection Pool

This package implements a PostgreSQL connection pool for the ByeBob application using `pgxpool` from the jackc/pgx library.

## Key Features

- Configurable connection pool settings
- Connection retry logic with exponential backoff
- Health check capabilities
- Global pool singleton for application-wide usage
- Helper methods for acquiring and releasing connections

## Usage Examples

### Basic Usage

```go
import (
    "context"
    "log"
    
    "github.com/gfurduy/byebob/config"
    "github.com/gfurduy/byebob/internal/repository"
)

func main() {
    // Load configuration
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize database pool with default settings
    db, err := repository.NewDBPool(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Use the pool
    // ...
}
```

### Using the Global Pool Singleton

```go
import (
    "context"
    "log"
    
    "github.com/gfurduy/byebob/config"
    "github.com/gfurduy/byebob/internal/repository"
)

func main() {
    // Load configuration
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize global pool (this only happens once)
    _, err = repository.InitGlobalDBPool(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize global database pool: %v", err)
    }
    defer repository.CloseGlobalDBPool()
    
    // Get the global pool from anywhere in your code
    pool := repository.GetGlobalDBPool()
    
    // Use the pool
    // ...
}
```

### Custom Pool Configuration

```go
import (
    "context"
    "log"
    "time"
    
    "github.com/gfurduy/byebob/config"
    "github.com/gfurduy/byebob/internal/repository"
)

func main() {
    // Load configuration
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Create custom pool configuration
    poolCfg := repository.DefaultPoolConfig()
    poolCfg.MaxConns = 20
    poolCfg.MinConns = 5
    poolCfg.MaxRetries = 10
    poolCfg.RetryDelay = 5 * time.Second
    
    // Initialize database pool with custom settings
    db, err := repository.NewDBPoolWithConfig(cfg, poolCfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Use the pool
    // ...
}
```

### Working with Connections

```go
import (
    "context"
    "log"
    
    "github.com/gfurduy/byebob/config"
    "github.com/gfurduy/byebob/internal/repository"
)

func main() {
    // Initialize the pool
    cfg, _ := config.NewConfig()
    db, _ := repository.NewDBPool(cfg)
    defer db.Close()
    
    // Method 1: Use WithAcquire to automatically acquire and release a connection
    ctx := context.Background()
    err := db.WithAcquire(ctx, func(conn *pgxpool.Conn) error {
        // Use the connection
        var result int
        return conn.QueryRow(ctx, "SELECT $1::int", 42).Scan(&result)
    })
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // Method 2: Manually acquire and release a connection
    conn, err := db.AcquireConn(ctx)
    if err != nil {
        log.Printf("Failed to acquire connection: %v", err)
        return
    }
    defer conn.Release()
    
    // Use the connection
    // ...
}
```

### Health Checks

```go
import (
    "context"
    "log"
    "time"
    
    "github.com/gfurduy/byebob/config"
    "github.com/gfurduy/byebob/internal/repository"
)

func main() {
    // Initialize the pool
    cfg, _ := config.NewConfig()
    db, _ := repository.NewDBPool(cfg)
    defer db.Close()
    
    // Run a health check
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.HealthCheck(ctx); err != nil {
        log.Printf("Database health check failed: %v", err)
    } else {
        log.Println("Database health check passed!")
    }
}
```

## Configuration Parameters

The `PoolConfig` struct allows you to configure the following parameters:

| Parameter | Description | Default Value |
|-----------|-------------|---------------|
| MaxConns | Maximum number of connections in the pool | 10 |
| MinConns | Minimum number of connections in the pool | 2 |
| MaxConnLifetime | Maximum lifetime of connections | 1 hour |
| MaxConnIdleTime | Maximum idle time of connections | 30 minutes |
| HealthCheckPeriod | Period between health checks | 1 minute |
| ConnectTimeout | Timeout for establishing connections | 5 seconds |
| MaxRetries | Maximum number of connection retries | 5 |
| RetryDelay | Initial delay between retries | 3 seconds |

## Best Practices

1. **Use the global singleton**: For most use cases, use the global singleton to avoid creating multiple pools.
2. **Always close the pool**: Make sure to call `Close()` when you're done with the pool to release resources.
3. **Use connection helpers**: Prefer `WithAcquire` over manually acquiring and releasing connections to avoid leaks.
4. **Set appropriate pool sizes**: Set `MaxConns` based on your database's capacity and your application's needs.
5. **Use timeouts**: Always use contexts with timeouts for database operations to avoid hanging.
6. **Handle retries at operation level**: The pool handles connection retries, but you may need to implement retries for specific operations.
7. **Regular health checks**: Implement periodic health checks in long-running applications.

## Monitoring

The `Stats()` method provides metrics about the pool that you can use for monitoring:

```go
stats := db.Stats()
log.Printf("Pool stats: %d acquired, %d constructed, %d total", 
    stats.AcquiredConns(), stats.ConstructedConns(), stats.TotalConns())
``` 