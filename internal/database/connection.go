package database

import (
	"github.com/tinwritescode/myapp/internal/dto/common"
	"github.com/tinwritescode/myapp/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(dsn string) error {
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		return common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to connect to database", err)
	}

	logger.Info("Database connected successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
