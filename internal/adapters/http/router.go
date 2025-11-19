package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fitness-tracker/internal/adapters/http/handlers"
	"fitness-tracker/internal/core/ports"
	"fitness-tracker/internal/config"
)

// SetupRouter initializes the HTTP router with all routes and middleware
func SetupRouter(
	authHandler *handlers.AuthHandler,
	authService ports.AuthService,
	cfg *config.Config,
) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Global middleware
	router.Use(corsMiddleware(cfg))
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())

	// Health check endpoint
	router.GET("/health", healthCheckHandler)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// TODO: Add other protected routes here
		// For now, only auth endpoints are enabled
	}

	return router
}

// corsMiddleware configures CORS settings
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "*"}, // Add your frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// requestIDMiddleware adds a unique request ID to each request
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("requestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// healthCheckHandler returns the health status of the API
func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "fitness-coach-api",
		"version": "1.0.0",
		"time":    time.Now().UTC(),
	})
}
