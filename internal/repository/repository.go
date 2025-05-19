package repository

import (
	"context"
	
	"github.com/gfurduy/byebob/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository is the interface for data access
type Repository interface {
	// Employee operations
	GetEmployees(ctx context.Context) ([]models.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (models.Employee, error)
	CreateEmployee(ctx context.Context, employee models.Employee) (string, error)
	UpdateEmployee(ctx context.Context, employee models.Employee) error
	DeleteEmployee(ctx context.Context, id string) error
}

// PostgresRepository implements Repository interface for PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

// Repository implementations will be added here as the project progresses 