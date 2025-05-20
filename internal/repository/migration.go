package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/gfurduy/byebob/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
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

	// Create the postgres driver
	driver, err := postgres.WithInstance(m.pool.Conn().Conn(), &postgres.Config{})
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
	// Implementation will be added
	return fmt.Errorf("not implemented yet")
}

// RollbackMigration rolls back the last migration
func (m *MigrationManager) RollbackMigration(migrationsPath string) error {
	// Create the postgres driver
	driver, err := postgres.WithInstance(m.pool.Conn().Conn(), &postgres.Config{})
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
	// Create the postgres driver
	driver, err := postgres.WithInstance(m.pool.Conn().Conn(), &postgres.Config{})
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