package main

import (
	"fmt"

	"github.com/tinwritescode/myapp/internal/config"
	"github.com/tinwritescode/myapp/internal/database"
	"github.com/tinwritescode/myapp/internal/models"
	"github.com/tinwritescode/myapp/internal/routes"
	"github.com/tinwritescode/myapp/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	if err := database.ConnectDB(dsn); err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.AutoMigrate(&models.User{}, &models.Account{}); err != nil {
		logger.Fatal("Failed to run migrations:", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := ":" + cfg.Server.Port
	logger.Infof("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
