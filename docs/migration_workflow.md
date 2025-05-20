# Database Migration Workflow

## Overview

This document explains how to use the migration system to manage database schema changes in the ByeBob application. The migration system uses golang-migrate to version-control and apply schema changes in a structured way.

## Directory Structure

Migrations are stored in the `migrations/postgres` directory, with each migration consisting of two files:
- `<sequence>_<description>.up.sql`: Contains the SQL statements to apply the migration
- `<sequence>_<description>.down.sql`: Contains the SQL statements to rollback the migration

## Migration Commands

The Makefile provides several commands to work with migrations:

### Setup Migrations

If you're starting from scratch, you can set up the migrations directory:

```bash
make setup-migrations
```

### Create a New Migration

To create a new migration:

```bash
make migrate-create name=migration_name
```

This will create two empty files in the `migrations/postgres` directory:
- `<timestamp>_migration_name.up.sql`
- `<timestamp>_migration_name.down.sql`

You should then edit these files to add your schema changes (in the .up.sql file) and rollback commands (in the .down.sql file).

### Apply Migrations

To apply all pending migrations:

```bash
make migrate-up
```

This command requires the `POSTGRESQL_URL` environment variable to be set, which should contain the PostgreSQL connection string.

### Rollback Last Migration

To rollback the most recently applied migration:

```bash
make migrate-down
```

## Migration Best Practices

1. **Transactional Migrations**: Always wrap your migrations in transactions (`BEGIN` and `COMMIT`) to ensure atomicity.

2. **Idempotent Migrations**: Use `IF EXISTS` and `IF NOT EXISTS` clauses to make migrations idempotent (can be run multiple times without error).

3. **Rollback Safety**: Ensure that down migrations properly undo the changes made in up migrations.

4. **Naming Conventions**: Use descriptive names for migrations, preferably with a sequence prefix.

5. **Testing**: Test migrations in a development environment before applying them to production.

6. **Documentation**: Add comments in migration files to explain complex changes.

## How It Works

The migration system:
1. Connects to the database using your PostgreSQL connection string
2. Creates a `schema_migrations` table if it doesn't exist
3. Checks which migrations have already been applied
4. Applies any pending migrations in order
5. Records each successful migration in the `schema_migrations` table

## Troubleshooting

### Migration Failed

If a migration fails:
1. Check the error message to understand what went wrong
2. Fix the issue in your migration file
3. Try running the migration again

### Dirty Database State

If the database is in a "dirty" state (a migration failed halfway):
1. Manually fix the database issues
2. Reset the migration version in the `schema_migrations` table
3. Run migrations again

## Implementation Details

The migration system is implemented in the `internal/repository/migration.go` file, which uses the golang-migrate library to handle migrations. 