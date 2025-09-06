package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tinwritescode/myapp/docs"
	"github.com/tinwritescode/myapp/internal/handlers"
)

func SetupRoutes(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", handlers.Ping)
		v1.GET("/users", handlers.GetUsers)
		v1.POST("/users", handlers.CreateUser)
	}

}
