package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tinwritescode/myapp/pkg/logger"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using environment variables or defaults")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "myapp"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDatabaseDSN returns the database connection string
// It prioritizes DATABASE_URL (used by Fly.io and Neon.db) over individual variables
func (c *Config) GetDatabaseDSN() string {
	// Check if DATABASE_URL is set (Fly.io or Neon.db provides this)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	// Fall back to individual variables for local development
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
	)
}
