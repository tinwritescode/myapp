package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/tinwritescode/myapp/internal/config"
	"github.com/tinwritescode/myapp/internal/database"
	"github.com/tinwritescode/myapp/internal/middleware"
	"github.com/tinwritescode/myapp/internal/models"
	"github.com/tinwritescode/myapp/internal/routes"
	"github.com/tinwritescode/myapp/internal/service"
	"github.com/tinwritescode/myapp/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize JWT secret
	middleware.SetJWTSecret(cfg.JWT.Secret)
	service.SetJWTSecret(cfg.JWT.Secret)

	// Connect to database
	dsn := cfg.GetDatabaseDSN()

	if err := database.ConnectDB(dsn); err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.AutoMigrate(&models.User{}, &models.Account{}, &models.URL{}, &models.RefreshToken{}); err != nil {
		logger.Fatal("Failed to run migrations:", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "https://myapp-frontend.fly.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := ":" + cfg.Server.Port
	logger.Infof("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
