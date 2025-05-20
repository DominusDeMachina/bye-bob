package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gfurduy/byebob/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	pool   *pgxpool.Pool
	config *config.Config
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(pool *pgxpool.Pool, cfg *config.Config) *MigrationManager {
	return &MigrationManager{
		pool:   pool,
		config: cfg,
	}
}

// RunMigrations applies all pending migrations
func (m *MigrationManager) RunMigrations(migrationsPath string) error {
	log.Println("Running database migrations from:", migrationsPath)

	// Create a sql.DB instance from the pgx pool
	db := stdlib.OpenDBFromPool(m.pool)
	defer db.Close()

	// Create the postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create the migrate instance
	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Run migrations
	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

// CreateMigration creates a new migration
func (m *MigrationManager) CreateMigration(name, migrationsPath string) error {
	log.Printf("Creating new migration '%s' in %s\n", name, migrationsPath)
	
	// Validate migration name
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}
	
	// Normalize migration name to use underscores instead of spaces
	normalizedName := strings.ReplaceAll(name, " ", "_")
	normalizedName = strings.ToLower(normalizedName)
	
	// Create up and down migration files
	timestamp := time.Now().Unix()
	upMigrationFileName := fmt.Sprintf("%s/%d_%s.up.sql", migrationsPath, timestamp, normalizedName)
	downMigrationFileName := fmt.Sprintf("%s/%d_%s.down.sql", migrationsPath, timestamp, normalizedName)
	
	// Write template content to the up migration file
	upContent := fmt.Sprintf(`-- Migration: %s (up)
-- Created at: %s

BEGIN;

-- Add your schema changes here

COMMIT;
`, name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upMigrationFileName, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}
	
	// Write template content to the down migration file
	downContent := fmt.Sprintf(`-- Migration: %s (down)
-- Created at: %s

BEGIN;

-- Add your rollback commands here

COMMIT;
`, name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(downMigrationFileName, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}
	
	log.Printf("Created migration files:\n- %s\n- %s\n", upMigrationFileName, downMigrationFileName)
	return nil
}

// RollbackMigration rolls back the last migration
func (m *MigrationManager) RollbackMigration(migrationsPath string) error {
	// Create a sql.DB instance from the pgx pool
	db := stdlib.OpenDBFromPool(m.pool)
	defer db.Close()

	// Create the postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create the migrate instance
	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Rollback one migration
	if err := migrator.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to roll back")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

// GetMigrationVersion gets the current migration version
func (m *MigrationManager) GetMigrationVersion(migrationsPath string) (uint, bool, error) {
	// Create a sql.DB instance from the pgx pool
	db := stdlib.OpenDBFromPool(m.pool)
	defer db.Close()

	// Create the postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create the migrate instance
	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}

	// Get current version
	version, dirty, err := migrator.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
} 