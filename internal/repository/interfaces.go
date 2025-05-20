package repository

import (
	"context"
	"time"
)

// Employee represents an employee record
type Employee struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"first_name"`
	MiddleName      string    `json:"middle_name,omitempty"`
	LastName        string    `json:"last_name"`
	DisplayName     string    `json:"display_name"`
	Email           string    `json:"email"`
	Address         string    `json:"address,omitempty"`
	PositionID      string    `json:"position_id"`
	DepartmentID    string    `json:"department_id"`
	SiteID          string    `json:"site_id"`
	ManagerID       string    `json:"manager_id,omitempty"`
	EmploymentType  string    `json:"employment_type"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date,omitempty"`
	Status          string    `json:"status"`
	ProfilePicture  string    `json:"profile_picture_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Position represents a job position
type Position struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Requirements string    `json:"requirements,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Department represents a department
type Department struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	LeadID      string    `json:"lead_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Site represents a physical location
type Site struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EmployeeRepository defines operations for working with employees
type EmployeeRepository interface {
	// Create a new employee
	Create(ctx context.Context, employee *Employee) (string, error)
	
	// Get an employee by ID
	GetByID(ctx context.Context, id string) (*Employee, error)
	
	// Update an employee
	Update(ctx context.Context, employee *Employee) error
	
	// Delete an employee
	Delete(ctx context.Context, id string) error
	
	// List employees with optional filters
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Employee, int64, error)
	
	// Get employees by manager ID
	GetByManager(ctx context.Context, managerID string) ([]*Employee, error)
	
	// Get employees by department ID
	GetByDepartment(ctx context.Context, departmentID string) ([]*Employee, error)
}

// PositionRepository defines operations for working with positions
type PositionRepository interface {
	Create(ctx context.Context, position *Position) (string, error)
	GetByID(ctx context.Context, id string) (*Position, error)
	Update(ctx context.Context, position *Position) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Position, int64, error)
}

// DepartmentRepository defines operations for working with departments
type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) (string, error)
	GetByID(ctx context.Context, id string) (*Department, error)
	Update(ctx context.Context, department *Department) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Department, int64, error)
}

// SiteRepository defines operations for working with sites
type SiteRepository interface {
	Create(ctx context.Context, site *Site) (string, error)
	GetByID(ctx context.Context, id string) (*Site, error)
	Update(ctx context.Context, site *Site) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Site, int64, error)
}

// RepositoryFactory defines the repository factory interface
type RepositoryFactory interface {
	Employees() EmployeeRepository
	Positions() PositionRepository
	Departments() DepartmentRepository
	Sites() SiteRepository
	
	// WithTransaction starts a new transaction and returns a RepositoryFactory that uses it
	WithTransaction(ctx context.Context) (RepositoryFactory, error)
	
	// Commit commits the current transaction
	Commit() error
	
	// Rollback rolls back the current transaction
	Rollback() error
} 