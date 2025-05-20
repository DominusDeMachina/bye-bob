package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Define queryer interface for common transaction and pool methods
type queryer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// PostgresFactory implements the RepositoryFactory interface for PostgreSQL
type PostgresFactory struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

// NewPostgresFactory creates a new PostgreSQL repository factory
func NewPostgresFactory(pool *pgxpool.Pool) *PostgresFactory {
	return &PostgresFactory{
		pool: pool,
	}
}

// Employees returns an EmployeeRepository
func (f *PostgresFactory) Employees() EmployeeRepository {
	return &PostgresEmployeeRepository{factory: f}
}

// Positions returns a PositionRepository
func (f *PostgresFactory) Positions() PositionRepository {
	return &PostgresPositionRepository{factory: f}
}

// Departments returns a DepartmentRepository
func (f *PostgresFactory) Departments() DepartmentRepository {
	return &PostgresDepartmentRepository{factory: f}
}

// Sites returns a SiteRepository
func (f *PostgresFactory) Sites() SiteRepository {
	return &PostgresSiteRepository{factory: f}
}

// WithTransaction starts a new transaction and returns a RepositoryFactory that uses it
func (f *PostgresFactory) WithTransaction(ctx context.Context) (RepositoryFactory, error) {
	if f.tx != nil {
		return nil, errors.New("transaction already started")
	}

	tx, err := f.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return &PostgresFactory{
		pool: f.pool,
		tx:   tx,
	}, nil
}

// Commit commits the current transaction
func (f *PostgresFactory) Commit() error {
	if f.tx == nil {
		return errors.New("no transaction to commit")
	}

	err := f.tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	f.tx = nil
	return nil
}

// Rollback rolls back the current transaction
func (f *PostgresFactory) Rollback() error {
	if f.tx == nil {
		return errors.New("no transaction to rollback")
	}

	err := f.tx.Rollback(context.Background())
	if err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	f.tx = nil
	return nil
}

// getQueryer returns the appropriate queryer (transaction or pool)
func (f *PostgresFactory) getQueryer() queryer {
	if f.tx != nil {
		return f.tx
	}
	return f.pool
}

// PostgresEmployeeRepository implements EmployeeRepository for PostgreSQL
type PostgresEmployeeRepository struct {
	factory *PostgresFactory
}

// Create creates a new employee
func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *Employee) (string, error) {
	query := `
		INSERT INTO employees (
			first_name, middle_name, last_name, display_name, email, 
			address, position_id, department_id, site_id, manager_id, 
			employment_type, start_date, end_date, status, profile_picture_url
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`

	var id string
	err := r.factory.getQueryer().QueryRow(ctx, query,
		employee.FirstName, employee.MiddleName, employee.LastName, employee.DisplayName,
		employee.Email, employee.Address, employee.PositionID, employee.DepartmentID,
		employee.SiteID, employee.ManagerID, employee.EmploymentType, employee.StartDate,
		employee.EndDate, employee.Status, employee.ProfilePicture,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create employee: %w", err)
	}

	return id, nil
}

// GetByID retrieves an employee by ID
func (r *PostgresEmployeeRepository) GetByID(ctx context.Context, id string) (*Employee, error) {
	query := `
		SELECT id, first_name, middle_name, last_name, display_name, email, 
			address, position_id, department_id, site_id, manager_id, 
			employment_type, start_date, end_date, status, profile_picture_url,
			created_at, updated_at
		FROM employees
		WHERE id = $1
	`

	var employee Employee
	var endDate sql.NullTime

	err := r.factory.getQueryer().QueryRow(ctx, query, id).Scan(
		&employee.ID, &employee.FirstName, &employee.MiddleName, &employee.LastName,
		&employee.DisplayName, &employee.Email, &employee.Address, &employee.PositionID,
		&employee.DepartmentID, &employee.SiteID, &employee.ManagerID, &employee.EmploymentType,
		&employee.StartDate, &endDate, &employee.Status, &employee.ProfilePicture,
		&employee.CreatedAt, &employee.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("employee not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	if endDate.Valid {
		employee.EndDate = endDate.Time
	}

	return &employee, nil
}

// Update updates an employee
func (r *PostgresEmployeeRepository) Update(ctx context.Context, employee *Employee) error {
	query := `
		UPDATE employees
		SET first_name = $1, middle_name = $2, last_name = $3, display_name = $4,
			email = $5, address = $6, position_id = $7, department_id = $8,
			site_id = $9, manager_id = $10, employment_type = $11, start_date = $12,
			end_date = $13, status = $14, profile_picture_url = $15, updated_at = NOW()
		WHERE id = $16
	`

	result, err := r.factory.getQueryer().Exec(ctx, query,
		employee.FirstName, employee.MiddleName, employee.LastName, employee.DisplayName,
		employee.Email, employee.Address, employee.PositionID, employee.DepartmentID,
		employee.SiteID, employee.ManagerID, employee.EmploymentType, employee.StartDate,
		employee.EndDate, employee.Status, employee.ProfilePicture, employee.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("employee not found: %s", employee.ID)
	}

	return nil
}

// Delete deletes an employee
func (r *PostgresEmployeeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.factory.getQueryer().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("employee not found: %s", id)
	}

	return nil
}

// List lists employees with optional filters
func (r *PostgresEmployeeRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Employee, int64, error) {
	// Base query
	query := `
		SELECT id, first_name, middle_name, last_name, display_name, email, 
			address, position_id, department_id, site_id, manager_id, 
			employment_type, start_date, end_date, status, profile_picture_url,
			created_at, updated_at
		FROM employees
	`

	// Where clause and parameters
	where := ""
	params := []interface{}{}
	paramIndex := 1

	if len(filters) > 0 {
		where = " WHERE "
		for key, value := range filters {
			if paramIndex > 1 {
				where += " AND "
			}
			where += fmt.Sprintf("%s = $%d", key, paramIndex)
			params = append(params, value)
			paramIndex++
		}
	}

	// Add pagination
	query += where + fmt.Sprintf(" ORDER BY last_name, first_name LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	params = append(params, limit, offset)

	// Count query
	countQuery := "SELECT COUNT(*) FROM employees" + where

	// Execute count query
	var total int64
	err := r.factory.getQueryer().QueryRow(ctx, countQuery, params[:len(params)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	// Execute main query
	rows, err := r.factory.getQueryer().Query(ctx, query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}
	defer rows.Close()

	employees := []*Employee{}
	for rows.Next() {
		var employee Employee
		var endDate sql.NullTime

		err := rows.Scan(
			&employee.ID, &employee.FirstName, &employee.MiddleName, &employee.LastName,
			&employee.DisplayName, &employee.Email, &employee.Address, &employee.PositionID,
			&employee.DepartmentID, &employee.SiteID, &employee.ManagerID, &employee.EmploymentType,
			&employee.StartDate, &endDate, &employee.Status, &employee.ProfilePicture,
			&employee.CreatedAt, &employee.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan employee: %w", err)
		}

		if endDate.Valid {
			employee.EndDate = endDate.Time
		}

		employees = append(employees, &employee)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating employee rows: %w", err)
	}

	return employees, total, nil
}

// GetByManager retrieves employees by manager ID
func (r *PostgresEmployeeRepository) GetByManager(ctx context.Context, managerID string) ([]*Employee, error) {
	query := `
		SELECT id, first_name, middle_name, last_name, display_name, email, 
			address, position_id, department_id, site_id, manager_id, 
			employment_type, start_date, end_date, status, profile_picture_url,
			created_at, updated_at
		FROM employees
		WHERE manager_id = $1
		ORDER BY last_name, first_name
	`

	rows, err := r.factory.getQueryer().Query(ctx, query, managerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by manager: %w", err)
	}
	defer rows.Close()

	employees := []*Employee{}
	for rows.Next() {
		var employee Employee
		var endDate sql.NullTime

		err := rows.Scan(
			&employee.ID, &employee.FirstName, &employee.MiddleName, &employee.LastName,
			&employee.DisplayName, &employee.Email, &employee.Address, &employee.PositionID,
			&employee.DepartmentID, &employee.SiteID, &employee.ManagerID, &employee.EmploymentType,
			&employee.StartDate, &endDate, &employee.Status, &employee.ProfilePicture,
			&employee.CreatedAt, &employee.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}

		if endDate.Valid {
			employee.EndDate = endDate.Time
		}

		employees = append(employees, &employee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employee rows: %w", err)
	}

	return employees, nil
}

// GetByDepartment retrieves employees by department ID
func (r *PostgresEmployeeRepository) GetByDepartment(ctx context.Context, departmentID string) ([]*Employee, error) {
	query := `
		SELECT id, first_name, middle_name, last_name, display_name, email, 
			address, position_id, department_id, site_id, manager_id, 
			employment_type, start_date, end_date, status, profile_picture_url,
			created_at, updated_at
		FROM employees
		WHERE department_id = $1
		ORDER BY last_name, first_name
	`

	rows, err := r.factory.getQueryer().Query(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by department: %w", err)
	}
	defer rows.Close()

	employees := []*Employee{}
	for rows.Next() {
		var employee Employee
		var endDate sql.NullTime

		err := rows.Scan(
			&employee.ID, &employee.FirstName, &employee.MiddleName, &employee.LastName,
			&employee.DisplayName, &employee.Email, &employee.Address, &employee.PositionID,
			&employee.DepartmentID, &employee.SiteID, &employee.ManagerID, &employee.EmploymentType,
			&employee.StartDate, &endDate, &employee.Status, &employee.ProfilePicture,
			&employee.CreatedAt, &employee.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}

		if endDate.Valid {
			employee.EndDate = endDate.Time
		}

		employees = append(employees, &employee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employee rows: %w", err)
	}

	return employees, nil
}

// PostgresPositionRepository implements PositionRepository for PostgreSQL
type PostgresPositionRepository struct {
	factory *PostgresFactory
}

// Create creates a new position
func (r *PostgresPositionRepository) Create(ctx context.Context, position *Position) (string, error) {
	query := `
		INSERT INTO positions (
			title, description, requirements
		) VALUES ($1, $2, $3)
		RETURNING id
	`

	var id string
	err := r.factory.getQueryer().QueryRow(ctx, query,
		position.Title, position.Description, position.Requirements,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create position: %w", err)
	}

	return id, nil
}

// GetByID retrieves a position by ID
func (r *PostgresPositionRepository) GetByID(ctx context.Context, id string) (*Position, error) {
	query := `
		SELECT id, title, description, requirements, created_at, updated_at
		FROM positions
		WHERE id = $1
	`

	var position Position
	err := r.factory.getQueryer().QueryRow(ctx, query, id).Scan(
		&position.ID, &position.Title, &position.Description, &position.Requirements,
		&position.CreatedAt, &position.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("position not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	return &position, nil
}

// Update updates a position
func (r *PostgresPositionRepository) Update(ctx context.Context, position *Position) error {
	query := `
		UPDATE positions
		SET title = $1, description = $2, requirements = $3, updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.factory.getQueryer().Exec(ctx, query,
		position.Title, position.Description, position.Requirements, position.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("position not found: %s", position.ID)
	}

	return nil
}

// Delete deletes a position
func (r *PostgresPositionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM positions WHERE id = $1`

	result, err := r.factory.getQueryer().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("position not found: %s", id)
	}

	return nil
}

// List lists positions
func (r *PostgresPositionRepository) List(ctx context.Context, limit, offset int) ([]*Position, int64, error) {
	// Count query
	countQuery := "SELECT COUNT(*) FROM positions"
	var total int64
	err := r.factory.getQueryer().QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count positions: %w", err)
	}

	// Main query
	query := `
		SELECT id, title, description, requirements, created_at, updated_at
		FROM positions
		ORDER BY title
		LIMIT $1 OFFSET $2
	`

	rows, err := r.factory.getQueryer().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list positions: %w", err)
	}
	defer rows.Close()

	positions := []*Position{}
	for rows.Next() {
		var position Position
		err := rows.Scan(
			&position.ID, &position.Title, &position.Description, &position.Requirements,
			&position.CreatedAt, &position.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan position: %w", err)
		}

		positions = append(positions, &position)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating position rows: %w", err)
	}

	return positions, total, nil
}

// PostgresDepartmentRepository implements DepartmentRepository for PostgreSQL
type PostgresDepartmentRepository struct {
	factory *PostgresFactory
}

// Create creates a new department
func (r *PostgresDepartmentRepository) Create(ctx context.Context, department *Department) (string, error) {
	query := `
		INSERT INTO departments (
			name, description, lead_id
		) VALUES ($1, $2, $3)
		RETURNING id
	`

	var id string
	err := r.factory.getQueryer().QueryRow(ctx, query,
		department.Name, department.Description, department.LeadID,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create department: %w", err)
	}

	return id, nil
}

// GetByID retrieves a department by ID
func (r *PostgresDepartmentRepository) GetByID(ctx context.Context, id string) (*Department, error) {
	query := `
		SELECT id, name, description, lead_id, created_at, updated_at
		FROM departments
		WHERE id = $1
	`

	var department Department
	err := r.factory.getQueryer().QueryRow(ctx, query, id).Scan(
		&department.ID, &department.Name, &department.Description, &department.LeadID,
		&department.CreatedAt, &department.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("department not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return &department, nil
}

// Update updates a department
func (r *PostgresDepartmentRepository) Update(ctx context.Context, department *Department) error {
	query := `
		UPDATE departments
		SET name = $1, description = $2, lead_id = $3, updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.factory.getQueryer().Exec(ctx, query,
		department.Name, department.Description, department.LeadID, department.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("department not found: %s", department.ID)
	}

	return nil
}

// Delete deletes a department
func (r *PostgresDepartmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM departments WHERE id = $1`

	result, err := r.factory.getQueryer().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("department not found: %s", id)
	}

	return nil
}

// List lists departments
func (r *PostgresDepartmentRepository) List(ctx context.Context, limit, offset int) ([]*Department, int64, error) {
	// Count query
	countQuery := "SELECT COUNT(*) FROM departments"
	var total int64
	err := r.factory.getQueryer().QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count departments: %w", err)
	}

	// Main query
	query := `
		SELECT id, name, description, lead_id, created_at, updated_at
		FROM departments
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.factory.getQueryer().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list departments: %w", err)
	}
	defer rows.Close()

	departments := []*Department{}
	for rows.Next() {
		var department Department
		err := rows.Scan(
			&department.ID, &department.Name, &department.Description, &department.LeadID,
			&department.CreatedAt, &department.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan department: %w", err)
		}

		departments = append(departments, &department)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating department rows: %w", err)
	}

	return departments, total, nil
}

// PostgresSiteRepository implements SiteRepository for PostgreSQL
type PostgresSiteRepository struct {
	factory *PostgresFactory
}

// Create creates a new site
func (r *PostgresSiteRepository) Create(ctx context.Context, site *Site) (string, error) {
	query := `
		INSERT INTO sites (
			name, city, address
		) VALUES ($1, $2, $3)
		RETURNING id
	`

	var id string
	err := r.factory.getQueryer().QueryRow(ctx, query,
		site.Name, site.City, site.Address,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create site: %w", err)
	}

	return id, nil
}

// GetByID retrieves a site by ID
func (r *PostgresSiteRepository) GetByID(ctx context.Context, id string) (*Site, error) {
	query := `
		SELECT id, name, city, address, created_at, updated_at
		FROM sites
		WHERE id = $1
	`

	var site Site
	err := r.factory.getQueryer().QueryRow(ctx, query, id).Scan(
		&site.ID, &site.Name, &site.City, &site.Address,
		&site.CreatedAt, &site.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("site not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	return &site, nil
}

// Update updates a site
func (r *PostgresSiteRepository) Update(ctx context.Context, site *Site) error {
	query := `
		UPDATE sites
		SET name = $1, city = $2, address = $3, updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.factory.getQueryer().Exec(ctx, query,
		site.Name, site.City, site.Address, site.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update site: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("site not found: %s", site.ID)
	}

	return nil
}

// Delete deletes a site
func (r *PostgresSiteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sites WHERE id = $1`

	result, err := r.factory.getQueryer().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("site not found: %s", id)
	}

	return nil
}

// List lists sites
func (r *PostgresSiteRepository) List(ctx context.Context, limit, offset int) ([]*Site, int64, error) {
	// Count query
	countQuery := "SELECT COUNT(*) FROM sites"
	var total int64
	err := r.factory.getQueryer().QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sites: %w", err)
	}

	// Main query
	query := `
		SELECT id, name, city, address, created_at, updated_at
		FROM sites
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.factory.getQueryer().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sites: %w", err)
	}
	defer rows.Close()

	sites := []*Site{}
	for rows.Next() {
		var site Site
		err := rows.Scan(
			&site.ID, &site.Name, &site.City, &site.Address,
			&site.CreatedAt, &site.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan site: %w", err)
		}

		sites = append(sites, &site)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating site rows: %w", err)
	}

	return sites, total, nil
} 