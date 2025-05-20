package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gfurduy/byebob/config"
	"github.com/gfurduy/byebob/internal/handlers"
	"github.com/gfurduy/byebob/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Version and BuildTime are set during build
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Print version info
	fmt.Printf("ByeBob %s (built at %s)\n", Version, BuildTime)

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Initialize database connection pool
	db, err := repository.NewDBPool(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// Initialize repository
	repo := repository.NewPostgresRepository(db.GetPool())
	
	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "ByeBob App",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	// Use global middlewares
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins[0],
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Static files
	app.Static("/static", "./static")

	// Setup routes
	handlers.SetupRoutes(app, repo)

	// Start the server in a goroutine
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	fmt.Printf("Server started on http://localhost:%s in %s mode\n", 
		cfg.Port, cfg.Environment)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	fmt.Println("Server gracefully stopped")
}

// customErrorHandler handles errors thrown by Fiber handlers
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code is 500
	statusCode := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		statusCode = e.Code
	}

	// Return JSON error response
	return c.Status(statusCode).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
