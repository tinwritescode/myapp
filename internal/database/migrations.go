package database

import (
	"log"

	"github.com/tinwritescode/myapp/internal/dto/common"
)

// AutoMigrate runs database migrations for all models
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return common.NewAppError(common.INTERNAL_SERVER_ERROR, "database connection not initialized", nil)
	}

	err := DB.AutoMigrate(models...)
	if err != nil {
		return common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to run migrations", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
