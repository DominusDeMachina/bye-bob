package repository

import (
	"github.com/gfurduy/byebob/internal/database"
)

// NewFactory creates a new repository factory using the database connection
func NewFactory(db *database.DB) RepositoryFactory {
	return NewPostgresFactory(db.Pool)
} 