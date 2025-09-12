package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tinwritescode/myapp/docs"
	"github.com/tinwritescode/myapp/internal/handlers"
	"github.com/tinwritescode/myapp/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Public routes (no authentication required)
	public := r.Group("/api/v1")
	{
		public.GET("/ping", handlers.Ping)
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/refresh", handlers.RefreshToken)
	}

	// Protected routes (authentication required)
	protected := r.Group("/api/v1").Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", handlers.GetUsers)

		// URL routes
		protected.POST("/urls", handlers.CreateURL)
		protected.GET("/urls", handlers.GetURLs)
		protected.GET("/urls/:id", handlers.GetURLByID)
		protected.PUT("/urls/:id", handlers.UpdateURL)
		protected.DELETE("/urls/:id", handlers.DeleteURL)
		protected.GET("/urls/:id/stats", handlers.GetURLStats)
	}

	// URL redirection route (outside API group for shorter URLs)
	r.GET("/:short_code", handlers.RedirectURL)
}
