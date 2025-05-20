package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gfurduy/byebob/internal/models"
	"github.com/jackc/pgx/v5"
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

// GetEmployees retrieves all employees from the database
func (r *PostgresRepository) GetEmployees(ctx context.Context) ([]models.Employee, error) {
	query := `
        SELECT id, created_at, updated_at, first_name, last_name, email, 
               position, department, manager_id, job_title, hire_date, 
               status, phone_number
        FROM employees
        ORDER BY last_name, first_name
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying employees: %w", err)
	}
	defer rows.Close()

	employees := []models.Employee{}
	for rows.Next() {
		var e models.Employee
		err := rows.Scan(
			&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.FirstName, &e.LastName,
			&e.Email, &e.Position, &e.Department, &e.ManagerID, &e.JobTitle,
			&e.HireDate, &e.Status, &e.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee row: %w", err)
		}
		employees = append(employees, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employee rows: %w", err)
	}

	return employees, nil
}

// GetEmployeeByID retrieves a single employee by ID
func (r *PostgresRepository) GetEmployeeByID(ctx context.Context, id string) (models.Employee, error) {
	query := `
        SELECT id, created_at, updated_at, first_name, last_name, email, 
               position, department, manager_id, job_title, hire_date, 
               status, phone_number
        FROM employees
        WHERE id = $1
    `

	var e models.Employee
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.FirstName, &e.LastName,
		&e.Email, &e.Position, &e.Department, &e.ManagerID, &e.JobTitle,
		&e.HireDate, &e.Status, &e.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Employee{}, fmt.Errorf("employee not found: %s", id)
		}
		return models.Employee{}, fmt.Errorf("error querying employee: %w", err)
	}

	return e, nil
}

// CreateEmployee creates a new employee in the database
func (r *PostgresRepository) CreateEmployee(ctx context.Context, e models.Employee) (string, error) {
	query := `
        INSERT INTO employees (
            first_name, last_name, email, position, department, 
            manager_id, job_title, hire_date, status, phone_number
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `

	var id string
	err := r.pool.QueryRow(ctx, query,
		e.FirstName, e.LastName, e.Email, e.Position, e.Department,
		e.ManagerID, e.JobTitle, e.HireDate, e.Status, e.PhoneNumber,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("error creating employee: %w", err)
	}

	return id, nil
}

// UpdateEmployee updates an existing employee in the database
func (r *PostgresRepository) UpdateEmployee(ctx context.Context, e models.Employee) error {
	query := `
        UPDATE employees
        SET first_name = $1, last_name = $2, email = $3, position = $4,
            department = $5, manager_id = $6, job_title = $7, hire_date = $8,
            status = $9, phone_number = $10, updated_at = NOW()
        WHERE id = $11
    `

	result, err := r.pool.Exec(ctx, query,
		e.FirstName, e.LastName, e.Email, e.Position, e.Department,
		e.ManagerID, e.JobTitle, e.HireDate, e.Status, e.PhoneNumber, e.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("employee not found: %s", e.ID)
	}

	return nil
}

// DeleteEmployee deletes an employee from the database
func (r *PostgresRepository) DeleteEmployee(ctx context.Context, id string) error {
	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("employee not found: %s", id)
	}

	return nil
}

// Repository implementations will be added here as the project progresses 