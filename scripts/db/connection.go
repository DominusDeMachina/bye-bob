package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gfurduy/byebob/config"
	"github.com/gfurduy/byebob/internal/database"
)

// TestDatabaseConnection tests the connection to a PostgreSQL database
func TestDatabaseConnection() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting database connection test...")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Display connection details
	fmt.Println("Database Connection Settings:")
	if cfg.RailwayDBURL != "" {
		fmt.Println("- Using Railway PostgreSQL database")
		fmt.Printf("- Railway DB URL: %s (masked)\n", maskConnectionString(cfg.RailwayDBURL))
	} else {
		fmt.Println("- Using direct PostgreSQL connection")
		fmt.Printf("- Host: %s\n", cfg.DBHost)
		fmt.Printf("- Port: %s\n", cfg.DBPort)
		fmt.Printf("- Database: %s\n", cfg.DBName)
		fmt.Printf("- User: %s\n", cfg.DBUser)
		fmt.Printf("- SSL Mode: %s\n", cfg.DBSSLMode)
	}

	fmt.Println("\nAttempting to connect to PostgreSQL...")
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify connection
	var version string
	err = db.Pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("\n✅ Connection successful!")
	fmt.Printf("PostgreSQL version: %s\n", version)

	// Check for database extensions
	var hasExtension bool
	err = db.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_graphql')").Scan(&hasExtension)
	if err != nil {
		fmt.Printf("Warning: Unable to check for database extensions: %v\n", err)
	} else if hasExtension {
		fmt.Println("Note: pg_graphql extension detected")
	}

	// Get table count
	var tableCount int
	err = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM pg_catalog.pg_tables 
		WHERE schemaname != 'pg_catalog' 
		AND schemaname != 'information_schema'
	`).Scan(&tableCount)
	if err != nil {
		fmt.Printf("Warning: Unable to get table count: %v\n", err)
	} else {
		fmt.Printf("Database contains %d user-defined tables\n", tableCount)
	}

	fmt.Println("\nTest completed successfully.")
}

// TestLocalConnection tests the connection to a local PostgreSQL database
func TestLocalConnection() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting local PostgreSQL connection test...")
	
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Force local connection settings
	cfg.RailwayDBURL = ""

	// Print connection details
	fmt.Println("Connection Settings:")
	fmt.Printf("- Host: %s\n", cfg.DBHost)
	fmt.Printf("- Port: %s\n", cfg.DBPort)
	fmt.Printf("- User: %s\n", cfg.DBUser)
	fmt.Printf("- Database: %s\n", cfg.DBName)
	fmt.Printf("- SSL Mode: %s\n", cfg.DBSSLMode)
	
	// Get the connection string
	connStr := cfg.PostgresConnectionString()
	fmt.Printf("- Connection String: %s\n", connStr)

	// Initialize database connection
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run a simple query to test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var version string
	err = db.Pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("\nConnection successful!")
	fmt.Printf("PostgreSQL version: %s\n", version)

	fmt.Println("\nTest completed successfully.")
}

// TestRailwayConnection tests the connection to a Railway PostgreSQL database
func TestRailwayConnection() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Railway database connection test...")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.RailwayDBURL == "" {
		log.Fatalf("Error: Railway Database URL must be configured in the environment variables.\n" +
			"Please ensure you have set the following in your .env file:\n" +
			"RAILWAY_DB_URL=postgresql://username:password@host.railway.app:port/database")
	}

	fmt.Println("Railway Connection Settings:")
	fmt.Printf("- Railway DB URL: %s (masked)\n", maskConnectionString(cfg.RailwayDBURL))
	
	fmt.Println("\nAttempting to connect to Railway PostgreSQL...")
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Railway database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var version string
	err = db.Pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("\n✅ Connection successful!")
	fmt.Printf("PostgreSQL version: %s\n", version)
	
	// Get database statistics
	var dbSize string
	err = db.Pool.QueryRow(ctx, "SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)
	if err != nil {
		log.Printf("Warning: Unable to get database size: %v", err)
	} else {
		fmt.Printf("Database size: %s\n", dbSize)
	}
	
	// Get table count
	var tableCount int
	err = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM pg_catalog.pg_tables 
		WHERE schemaname != 'pg_catalog' 
		AND schemaname != 'information_schema'
	`).Scan(&tableCount)
	if err != nil {
		fmt.Printf("Warning: Unable to get table count: %v\n", err)
	} else {
		fmt.Printf("Database contains %d user-defined tables\n", tableCount)
	}

	fmt.Println("\nTest completed successfully.")
}

// maskConnectionString masks the password in a connection string for security
func maskConnectionString(connStr string) string {
	// Basic masking implementation
	if connStr == "" {
		return ""
	}

	// Find the password part in the connection string
	// Format: postgres://username:password@host:port/database
	parts := strings.Split(connStr, "@")
	if len(parts) < 2 {
		// If no @ symbol, return a completely masked string
		return "postgres://****:****@****:****/*****"
	}

	credentials := strings.Split(parts[0], ":")
	if len(credentials) < 3 {
		// If no password part, return a masked string
		return "postgres://****:****@" + parts[1]
	}

	// Mask the password but keep the protocol, username, and the rest
	protocol := credentials[0]
	username := credentials[1]
	maskedConnStr := fmt.Sprintf("%s:%s:****@%s", protocol, username, parts[1])
	
	return maskedConnStr
} 