package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server config
	Port          string
	Environment   string
	AllowedOrigins []string

	// Database config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	
	// Railway config
	RailwayDBURL string

	// Clerk auth config
	ClerkSecretKey string
	ClerkPubKey    string
}

// NewConfig loads configuration from environment variables
func NewConfig() (*Config, error) {
	// Determine environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Load env file based on environment
	envFile := ".env"
	if env != "development" {
		envFile = fmt.Sprintf(".env.%s", env)
	}

	// Load .env file if it exists
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("error loading %s file: %w", envFile, err)
		}
	}

	// Create config
	cfg := &Config{
		// Server config
		Port:          getEnv("PORT", "3000"),
		Environment:   env,
		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ","),

		// Database config
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "byebob"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		
		// Railway config
		RailwayDBURL: getEnv("RAILWAY_DB_URL", ""),

		// Clerk auth config
		ClerkSecretKey: getEnv("CLERK_SECRET_KEY", ""),
		ClerkPubKey:    getEnv("CLERK_PUB_KEY", ""),
	}

	return cfg, nil
}

// PostgresConnectionString returns the PostgreSQL connection string
func (c *Config) PostgresConnectionString() string {
	// If Railway DB URL is provided, use it
	if c.RailwayDBURL != "" {
		return c.RailwayDBURL
	}
	
	// Fallback to regular PostgreSQL connection
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

// Helper functions for retrieving environment variables
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetProjectRoot returns the absolute path to the project root
func GetProjectRoot() (string, error) {
	// Try to find the project root by looking for a .git directory or go.mod file
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	// If we can't find a project root, just return the current directory
	return os.Getwd()
}

// IsDevelopment returns true if in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
} 