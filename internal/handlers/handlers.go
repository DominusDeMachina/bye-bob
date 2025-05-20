package handlers

import (
	"github.com/gfurduy/byebob/config"
	"github.com/gfurduy/byebob/internal/repository"
	"github.com/gfurduy/byebob/internal/templates"
	"github.com/gofiber/fiber/v2"
)

// Handler manages the application's HTTP handlers
type Handler struct {
	repo repository.Repository
}

// NewHandler creates a new handler with the given repository
func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, repo repository.Repository) {
	// Create a handler with the repository
	h := NewHandler(repo)

	// Web routes (HTML)
	app.Get("/", HomeHandler)

	// API v1 routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health check
	v1.Get("/health", HealthCheck)

	// Employee routes
	employees := v1.Group("/employees")
	employees.Get("/", h.GetEmployees)
	employees.Get("/:id", h.GetEmployee)
	// Add more employee routes as needed
}

// HomeHandler renders the home page
func HomeHandler(c *fiber.Ctx) error {
	return templates.Home().Render(c.Context(), c.Response().BodyWriter())
}

// HealthCheck handler for health endpoint
func HealthCheck(c *fiber.Ctx) error {
	cfg, _ := config.NewConfig()
	return c.JSON(fiber.Map{
		"status":      "ok",
		"message":     "ByeBob API is healthy",
		"environment": cfg.Environment,
		"version":     "dev", // This would come from the version in main
	})
}

// GetEmployees returns a list of employees
func (h *Handler) GetEmployees(c *fiber.Ctx) error {
	employees, err := h.repo.GetEmployees(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error fetching employees",
		})
	}
	
	return c.JSON(fiber.Map{
		"data": employees,
	})
}

// GetEmployee returns a single employee by ID
func (h *Handler) GetEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	
	employee, err := h.repo.GetEmployeeByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error fetching employee",
		})
	}
	
	return c.JSON(fiber.Map{
		"data": employee,
	})
} 