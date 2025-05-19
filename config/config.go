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

	// Try to load .env file, but don't error if it doesn't exist
	_ = godotenv.Load(envFile)

	// Create config from environment variables
	config := &Config{
		// Server config
		Port:          getEnv("PORT", "3000"),
		Environment:   env,
		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ","),

		// Database config
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "byebob"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Clerk auth config
		ClerkSecretKey: getEnv("CLERK_SECRET_KEY", ""),
		ClerkPubKey:    getEnv("CLERK_PUB_KEY", ""),
	}

	return config, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// IsDevelopment returns true if in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt gets an environment variable as int or returns a default value
func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvBool gets an environment variable as bool or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetRootDir returns the absolute path to the project root directory
func GetRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// Traverse up until we find go.mod
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// We've reached the filesystem root without finding go.mod
			break
		}
		dir = parentDir
	}
	
	return cwd, nil
} 