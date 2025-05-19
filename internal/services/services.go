package services

import (
	"context"
	
	"github.com/gfurduy/byebob/internal/models"
	"github.com/gfurduy/byebob/internal/repository"
)

// EmployeeService handles employee business logic
type EmployeeService struct {
	repo repository.Repository
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(repo repository.Repository) *EmployeeService {
	return &EmployeeService{
		repo: repo,
	}
}

// GetEmployees retrieves all employees
func (s *EmployeeService) GetEmployees(ctx context.Context) ([]models.Employee, error) {
	return s.repo.GetEmployees(ctx)
}

// GetEmployeeByID retrieves an employee by ID
func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id string) (models.Employee, error) {
	return s.repo.GetEmployeeByID(ctx, id)
}

// Other service implementations will be added here as the project progresses 