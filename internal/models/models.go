package models

import (
	"time"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Employee represents an employee in the system
type Employee struct {
	BaseModel
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Position    string `json:"position"`
	Department  string `json:"department"`
	ManagerID   string `json:"manager_id,omitempty"`
	JobTitle    string `json:"job_title"`
	HireDate    string `json:"hire_date"`
	Status      string `json:"status"` // active, inactive, terminated, etc.
	PhoneNumber string `json:"phone_number,omitempty"`
	// Additional fields will be added as needed
}

// Other models will be added here as the project progresses 