package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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

	// Verify we have Supabase configuration
	if cfg.SupabaseURL == "" || cfg.SupabaseService == "" {
		log.Fatalf("Error: Supabase URL and Service Key must be configured in the environment variables.\n"+
			"Please ensure you have set the following in your .env file:\n"+
			"SUPABASE_URL=https://your-project-id.supabase.co\n"+
			"SUPABASE_SERVICE_KEY=your-service-role-key\n\n"+
			"Current settings:\n"+
			"SUPABASE_URL=%s\n"+
			"SUPABASE_SERVICE_KEY=%s\n",
			cfg.SupabaseURL, 
			maskKey(cfg.SupabaseService))
	}

	// Print connection details (without sensitive info)
	fmt.Println("Supabase Connection Settings:")
	fmt.Printf("- Supabase URL: %s\n", cfg.SupabaseURL)
	fmt.Printf("- API Key: %s\n", maskKey(cfg.SupabaseAPIKey))
	fmt.Printf("- Anon Key: %s\n", maskKey(cfg.SupabaseAnon))
	fmt.Printf("- Service Key: %s\n", maskKey(cfg.SupabaseService))
	
	// Extract project ID from URL for display purposes
	projectID := extractProjectID(cfg.SupabaseURL)
	if projectID != "" {
		fmt.Printf("- Project ID: %s\n", projectID)
	}
	
	// Construct database host from Supabase URL
	dbHost := ""
	if projectID != "" {
		dbHost = fmt.Sprintf("db.%s.supabase.co", projectID)
		fmt.Printf("- DB Host: %s\n", dbHost)
	}

	// Initialize database connection
	fmt.Println("\nAttempting to connect to Supabase PostgreSQL...")
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Supabase: %v", err)
	}
	defer db.Close()

	// Run a simple query to test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var version string
	err = db.Pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("\n✅ Connection successful!")
	fmt.Printf("PostgreSQL version: %s\n", version)

	// Check for Supabase extensions
	var hasGraphQLExt bool
	err = db.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_graphql')").Scan(&hasGraphQLExt)
	if err != nil {
		fmt.Printf("Warning: Unable to check for Supabase extensions: %v\n", err)
	} else if hasGraphQLExt {
		fmt.Println("✅ Verified Supabase database (pg_graphql extension detected)")
	}

	fmt.Println("\nTest completed successfully.")
	os.Exit(0)
}

// maskKey returns a masked version of a key for display purposes
func maskKey(key string) string {
	if len(key) < 8 {
		return "[not set]"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

// extractProjectID extracts the project ID from a Supabase URL
func extractProjectID(url string) string {
	// URL format: https://project-id.supabase.co
	parts := strings.Split(url, ".")
	if len(parts) < 2 {
		return ""
	}
	
	hostPart := parts[0]
	return strings.TrimPrefix(hostPart, "https://")
} 