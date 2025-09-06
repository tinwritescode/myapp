package database

import (
	"fmt"
	"log"
)

// AutoMigrate runs database migrations for all models
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	err := DB.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

