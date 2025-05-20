package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gfurduy/byebob/config"
	"github.com/gfurduy/byebob/internal/database"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Supabase connection test...")

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print connection details (without sensitive info)
	fmt.Println("Connection Settings:")
	fmt.Printf("- Host: %s\n", cfg.DBHost)
	fmt.Printf("- Port: %s\n", cfg.DBPort)
	fmt.Printf("- User: %s\n", cfg.DBUser)
	fmt.Printf("- Database: %s\n", cfg.DBName)
	fmt.Printf("- SSL Mode: %s\n", cfg.DBSSLMode)
	
	if cfg.SupabaseURL != "" {
		fmt.Printf("- Supabase URL: %s\n", cfg.SupabaseURL)
		fmt.Println("- Supabase API Key: [REDACTED]")
	} else {
		fmt.Println("Supabase URL not configured. Using direct PostgreSQL connection.")
	}

	// Initialize database connection
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
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

	// Check if we're running in Supabase
	var isSupabase bool
	err = db.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_graphql')").Scan(&isSupabase)
	if err != nil {
		log.Printf("Warning: Failed to check for Supabase extensions: %v", err)
	} else if isSupabase {
		fmt.Println("This appears to be a Supabase database (pg_graphql extension detected).")
	}

	fmt.Println("\nTest completed successfully.")
	os.Exit(0)
} 