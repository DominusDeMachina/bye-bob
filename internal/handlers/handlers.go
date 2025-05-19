package handlers

import (
	"github.com/gfurduy/byebob/config"
	"github.com/gfurduy/byebob/internal/templates"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// Web routes (HTML)
	app.Get("/", HomeHandler)

	// API v1 routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health check
	v1.Get("/health", HealthCheck)

	// Employee routes
	employees := v1.Group("/employees")
	employees.Get("/", GetEmployees)
	employees.Get("/:id", GetEmployee)
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
func GetEmployees(c *fiber.Ctx) error {
	// This would normally fetch from the database
	return c.JSON(fiber.Map{
		"message": "This endpoint will return a list of employees",
	})
}

// GetEmployee returns a single employee by ID
func GetEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	// This would normally fetch from the database
	return c.JSON(fiber.Map{
		"message": "This endpoint will return an employee with ID: " + id,
	})
} 