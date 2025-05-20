# ByeBob - Enterprise HR Management System

ByeBob is a comprehensive HR management system built with Go, designed to support 10,000+ employees with features for administration, recruitment, and employee self-service.

## Features

- Centralized employee data management
- Performance appraisal system with customizable templates
- Goal management with progress tracking
- Employee portal with profile management

## Technical Stack

- **Backend**: Go with Fiber web framework
- **Templating**: Templ for type-safe HTML generation
- **Frontend**: HTMX + TailwindCSS
- **Database**: PostgreSQL on Railway.com
- **Authentication**: Clerk
- **Deployment**: Docker, Railway.com, or VPS

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL
- Docker (optional)

### Setup

1. Clone the repository
   ```bash
   git clone https://github.com/gfurduy/byebob.git
   cd byebob
   ```

2. Copy the environment file
   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your credentials

4. Initialize the database
   ```bash
   make migrate-up
   ```

5. Run the application
   ```bash
   make dev
   ```

## Database Configuration

ByeBob uses Railway.com PostgreSQL for its database. The application supports the following connection methods:

1. **Railway Database URL** (recommended):
   ```
   RAILWAY_DB_URL=postgresql://postgres:password@containers-us-west-xxx.railway.app:5432/railway
   ```

2. **Individual connection parameters**:
   ```
   DB_HOST=containers-us-west-xxx.railway.app
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your-password
   DB_NAME=railway
   DB_SSLMODE=require
   ```

### Database Migration from Supabase

If you're migrating from Supabase to Railway.com, follow the steps in [docs/supabase_to_railway_migration.md](docs/supabase_to_railway_migration.md).

## Development

### Available Commands

- `make dev` - Run the application with hot-reloading
- `make build` - Build the application
- `make run` - Run the application
- `make docker-dev` - Run the development environment with Docker Compose
- `make test` - Run tests
- `make lint` - Run linting
- `make test-railway` - Test Railway.com database connection

### Database Migrations

- `make migrate-up` - Apply all pending migrations
- `make migrate-down` - Roll back the last migration
- `make migrate-create name=migration_name` - Create a new migration

## Documentation

- [Database Setup](docs/database_setup.md)
- [Railway.com Setup](docs/railway_setup.md)
- [Migration from Supabase to Railway](docs/supabase_to_railway_migration.md)
- [Database Security](docs/database_security.md)
- [Migration Workflow](docs/migration_workflow.md)

## License

This project is licensed under the MIT License. 