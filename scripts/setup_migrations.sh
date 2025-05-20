#!/bin/bash

# Setup script for initializing the migration system

# Create necessary directories
mkdir -p migrations/postgres

# Check if golang-migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "golang-migrate not found. Installing..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    
    if [ $? -ne 0 ]; then
        echo "Failed to install golang-migrate. Please install it manually."
        exit 1
    fi
    
    echo "golang-migrate installed successfully."
else
    echo "golang-migrate is already installed."
fi

# Create initial migration
echo "Creating initial migration..."
migrate create -ext sql -dir migrations/postgres -seq initial_schema

if [ $? -ne 0 ]; then
    echo "Failed to create initial migration."
    exit 1
fi

# Add up/down templates to the migration files
LATEST_MIGRATION=$(ls -1 migrations/postgres/*_initial_schema.up.sql | head -n 1)
DOWN_MIGRATION=${LATEST_MIGRATION/.up./.down.}

if [ -f "$LATEST_MIGRATION" ]; then
    # Add comments to up migration
    cat > "$LATEST_MIGRATION" << EOF
-- Initial schema migration for ByeBob
-- 
-- This migration creates the core tables for the application
-- including employees, positions, departments, and sites.

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Add your schema creation SQL here

EOF
    
    # Add comments to down migration
    cat > "$DOWN_MIGRATION" << EOF
-- Rollback for initial schema migration
-- 
-- This migration drops all tables created in the up migration

-- Add your schema rollback SQL here

EOF
    
    echo "Migration templates created successfully at:"
    echo "  $LATEST_MIGRATION"
    echo "  $DOWN_MIGRATION"
    echo ""
    echo "Next steps:"
    echo "1. Edit the migration files to add your schema SQL"
    echo "2. Set the POSTGRESQL_URL environment variable"
    echo "3. Run 'make migrate-up' to apply the migration"
fi

# Create migrations directory README
cat > migrations/README.md << EOF
# Database Migrations

This directory contains database migrations for the ByeBob application.

## Directory Structure

- \`postgres/\` - PostgreSQL migrations using the golang-migrate format

## Migration Naming Convention

Migrations follow the format:
\`000001_description.up.sql\` and \`000001_description.down.sql\`

## Running Migrations

Use the following commands from the project root:

\`\`\`bash
# Apply all pending migrations
make migrate-up

# Roll back the last migration
make migrate-down

# Create a new migration
make migrate-create name=migration_name
\`\`\`

## Migration Best Practices

1. Always create both up and down migrations
2. Keep migrations small and focused
3. Test both directions (up and down)
4. Never edit an existing migration that has been applied
5. Include comments explaining complex changes
6. Use transactions for data consistency
7. Consider backward compatibility
EOF

echo "Setup complete! Migration system initialized." 