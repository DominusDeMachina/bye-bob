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
	log.Println("Starting local PostgreSQL connection test...")

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Force local connection by clearing Supabase settings
	cfg.SupabaseURL = ""
	cfg.SupabaseAPIKey = ""
	cfg.SupabaseAnon = ""
	cfg.SupabaseService = ""

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

	fmt.Println("\nTest completed successfully.")
	os.Exit(0)
} 