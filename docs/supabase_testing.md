# Testing Supabase Connection

This document explains how to test your Supabase connection after setting up the project.

## Prerequisites

Before testing, make sure you have:

1. Created a Supabase project (see [Supabase Setup Guide](supabase_setup.md))
2. Collected your Supabase credentials:
   - **Project URL**: `https://your-project-id.supabase.co`  
   - **Service Role Key**: The `service_role` key from the API section

## Setting Environment Variables

1. Create a `.env` file in the project root (if it doesn't exist already)
2. Add the following Supabase-specific variables:

```
# Supabase
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_API_KEY=your-api-key
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-role-key
```

Replace the placeholder values with your actual Supabase credentials.

## Testing the Connection

### Option 1: Run the Test Command

The simplest way to test your Supabase connection is to use the provided Makefile target:

```bash
make test-supabase
```

This command will:
1. Create and run a temporary test script
2. Attempt to connect to your Supabase PostgreSQL instance
3. Run a simple query to verify the connection
4. Check for Supabase-specific database extensions

### Option 2: Manual Test

If you prefer to test manually, you can:

1. Create a temporary Go file:

```bash
mkdir -p tmp
cd tmp
```

2. Create a file named `test.go` with the following content:

```go
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
		log.Fatalf("Supabase URL and Service Key must be configured")
	}

	// Initialize database connection
	fmt.Println("Attempting to connect to Supabase PostgreSQL...")
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

	fmt.Println("Connection successful!")
	fmt.Printf("PostgreSQL version: %s\n", version)
	fmt.Println("Test completed successfully.")
}
```

3. Run the test:

```bash
go run test.go
```

## Troubleshooting

If you encounter connection issues:

1. **Authentication Errors**:
   - Double-check your `SUPABASE_SERVICE_KEY` - you need the `service_role` key, not the `anon` key
   - Ensure your key has not expired and has the necessary permissions

2. **Connectivity Issues**:
   - Verify your network can reach Supabase (no firewall blocking access)
   - Try from a different network if possible

3. **Database Issues**:
   - Check if your database is paused (free tier Supabase databases pause after inactivity)
   - Visit your Supabase dashboard to check the status of your database
   - Use the SQL Editor in Supabase Studio to check if direct queries work

## Using the Web Interface

Sometimes the easiest way to test is directly in the Supabase web interface:

1. Log in to your Supabase account at [supabase.com](https://supabase.com)
2. Go to your project
3. Click on "SQL Editor" in the left sidebar
4. Run a simple query like `SELECT version();`
5. If this works, but your Go application can't connect, it's likely an authentication or networking issue 