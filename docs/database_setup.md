# Database Setup Guide

This document provides instructions for setting up the PostgreSQL database with Railway and managing migrations for the ByeBob application.

## Railway Setup

For detailed Railway setup instructions, see [Railway Setup Guide](railway_setup.md).

## Environment Configuration

1. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Update the database connection settings in your `.env` file:
   ```
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your-password
   DB_NAME=byebob
   DB_SSLMODE=disable
   
   # Railway
   RAILWAY_DB_URL=postgresql://username:password@host.railway.app:port/database
   ```

3. For migration commands, set the PostgreSQL connection URL:
   ```bash
   # For local development
   export POSTGRESQL_URL="postgres://postgres:your-password@localhost:5432/byebob?sslmode=disable"
   
   # For Railway
   export POSTGRESQL_URL="postgresql://username:password@host.railway.app:port/database"
   ```

## Migration System Setup

1. Run the migration setup script:
   ```bash
   ./scripts/setup_migrations.sh
   ```

   This script will:
   - Create the migrations directory structure
   - Install golang-migrate if not already installed
   - Create initial migration templates
   - Generate a README with migration best practices

2. Alternatively, you can set up the migrations manually:
   ```bash
   # Install golang-migrate
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # Create migrations directory
   make setup-migrations
   
   # Create initial migration
   make migrate-create name=initial_schema
   ```

## Creating Migrations

1. Create a new migration:
   ```bash
   make migrate-create name=create_employees_table
   ```

2. Edit the generated migration files in `migrations/postgres`:
   - `000001_create_employees_table.up.sql` - Schema changes to apply
   - `000001_create_employees_table.down.sql` - Rollback instructions

## Running Migrations

1. Apply migrations:
   ```bash
   make migrate-up
   ```

2. Roll back the last migration:
   ```bash
   make migrate-down
   ```

## Testing Database Connection

Test your database connection:
```bash
make test-db
```

This will:
- Attempt to connect to the database using the configured settings
- Run a simple query to verify the connection
- Display the PostgreSQL version and connection details

## Implementing Migrations

When implementing migrations, follow these guidelines:

1. Use transactions for data consistency:
   ```sql
   -- In up.sql
   BEGIN;
   -- Your schema changes here
   COMMIT;
   
   -- In down.sql
   BEGIN;
   -- Your rollback changes here
   COMMIT;
   ```

2. Always include both `up` and `down` migrations

3. Use explicit naming conventions for tables, columns, and constraints:
   ```sql
   CREATE TABLE employees (
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     first_name VARCHAR(100) NOT NULL,
     -- other columns
     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
     CONSTRAINT fk_department FOREIGN KEY (department_id) REFERENCES departments(id)
   );
   ```

4. Add indexes for frequently queried columns:
   ```sql
   CREATE INDEX idx_employees_email ON employees(email);
   ```

5. Implement audit triggers for tracking changes:
   ```sql
   CREATE TRIGGER employees_audit
   AFTER INSERT OR UPDATE OR DELETE ON employees
   FOR EACH ROW EXECUTE PROCEDURE audit_log();
   ```

## Connection Pooling

The application uses `pgxpool` for connection pooling. The pool is configured in `internal/database/db.go` with the following settings:

- Max connections: 10
- Min connections: 2
- Max connection lifetime: 1 hour
- Max connection idle time: 30 minutes
- Health check period: 1 minute

Adjust these values in production based on your expected load.

## Database Access Pattern

The application follows the repository pattern for database access:

1. Interface definitions in `internal/repository/interfaces.go`
2. PostgreSQL implementations in `internal/repository/postgres/`
3. Transaction support through the repository factory

This pattern provides clean separation and makes unit testing easier with mock implementations. 