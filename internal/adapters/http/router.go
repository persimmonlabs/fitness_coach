package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"fitness-tracker/internal/adapters/http/handlers"
	"fitness-tracker/internal/adapters/http/middleware"
	"fitness-tracker/internal/config"
)

// SetupRouter initializes the HTTP router with all routes and middleware
func SetupRouter(
	authHandler *handlers.AuthHandler,
	mealHandler *handlers.MealHandler,
	foodHandler *handlers.FoodHandler,
	activityHandler *handlers.ActivityHandler,
	workoutHandler *handlers.WorkoutHandler,
	exerciseHandler *handlers.ExerciseHandler,
	metricHandler *handlers.MetricHandler,
	goalHandler *handlers.GoalHandler,
	chatHandler *handlers.ChatHandler,
	summaryHandler *handlers.SummaryHandler,
	config *config.Config,
	logger *zap.Logger,
) *gin.Engine {
	// Set Gin mode based on environment
	if config.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Global middleware
	router.Use(corsMiddleware(config))
	router.Use(loggerMiddleware(logger))
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

		// Protected routes (authentication required)
		authMiddleware := middleware.AuthJWT(config.JWT.Secret, logger)

		// Meal routes
		meals := v1.Group("/meals")
		meals.Use(authMiddleware)
		{
			meals.POST("", mealHandler.Create)
			meals.GET("", mealHandler.List)
			meals.GET("/:id", mealHandler.GetByID)
			meals.PUT("/:id", mealHandler.Update)
			meals.DELETE("/:id", mealHandler.Delete)
		}

		// Food routes
		foods := v1.Group("/foods")
		foods.Use(authMiddleware)
		{
			foods.POST("", foodHandler.Create)
			foods.GET("", foodHandler.List)
			foods.GET("/:id", foodHandler.GetByID)
			foods.PUT("/:id", foodHandler.Update)
			foods.DELETE("/:id", foodHandler.Delete)
		}

		// Activity routes
		activities := v1.Group("/activities")
		activities.Use(authMiddleware)
		{
			activities.POST("", activityHandler.Create)
			activities.GET("", activityHandler.List)
			activities.GET("/:id", activityHandler.GetByID)
			activities.PUT("/:id", activityHandler.Update)
			activities.DELETE("/:id", activityHandler.Delete)
		}

		// Workout routes
		workouts := v1.Group("/workouts")
		workouts.Use(authMiddleware)
		{
			workouts.POST("", workoutHandler.Create)
			workouts.GET("", workoutHandler.List)
			workouts.GET("/:id", workoutHandler.GetByID)
			workouts.PUT("/:id", workoutHandler.Update)
			workouts.DELETE("/:id", workoutHandler.Delete)
		}

		// Exercise routes
		exercises := v1.Group("/exercises")
		exercises.Use(authMiddleware)
		{
			exercises.POST("", exerciseHandler.Create)
			exercises.GET("", exerciseHandler.List)
			exercises.GET("/:id", exerciseHandler.GetByID)
			exercises.PUT("/:id", exerciseHandler.Update)
			exercises.DELETE("/:id", exerciseHandler.Delete)
		}

		// Metric routes
		metrics := v1.Group("/metrics")
		metrics.Use(authMiddleware)
		{
			metrics.POST("", metricHandler.Create)
			metrics.GET("", metricHandler.List)
			metrics.GET("/:id", metricHandler.GetByID)
			metrics.PUT("/:id", metricHandler.Update)
			metrics.DELETE("/:id", metricHandler.Delete)
		}

		// Goal routes
		goals := v1.Group("/goals")
		goals.Use(authMiddleware)
		{
			goals.POST("", goalHandler.Create)
			goals.GET("", goalHandler.List)
			goals.GET("/:id", goalHandler.GetByID)
			goals.PUT("/:id", goalHandler.Update)
			goals.DELETE("/:id", goalHandler.Delete)
		}

		// Chat routes
		chat := v1.Group("/chat")
		chat.Use(authMiddleware)
		{
			chat.POST("/message", chatHandler.SendMessage)
			chat.GET("/history", chatHandler.GetHistory)
		}

		// Onboarding routes
		onboarding := v1.Group("/onboarding")
		onboarding.Use(authMiddleware)
		{
			onboarding.POST("/complete", authHandler.CompleteOnboarding)
		}

		// Summary route
		summary := v1.Group("/summary")
		summary.Use(authMiddleware)
		{
			summary.GET("/daily", summaryHandler.GetDailySummary)
		}
	}

	return router
}

// corsMiddleware configures CORS settings
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // Add your frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// loggerMiddleware logs HTTP requests
func loggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP request",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("error", errorMessage),
		)
	}
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
