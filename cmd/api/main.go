package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fitness-tracker/internal/adapters/http/handlers"
	httpAdapter "fitness-tracker/internal/adapters/http"
	"fitness-tracker/internal/adapters/repositories/postgres"
	"fitness-tracker/internal/config"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/services"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.Server.Environment == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	logger.Info("Starting Fitness Tracker API",
		zap.String("env", cfg.Server.Environment),
		zap.Int("port", cfg.Server.Port),
	)

	// Initialize database
	db, err := config.InitDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Food{},
		&domain.ServingUnit{},
		&domain.FoodServingConversion{},
		&domain.Meal{},
		&domain.MealFoodItem{},
		&domain.Activity{},
		&domain.Exercise{},
		&domain.Workout{},
		&domain.WorkoutExercise{},
		&domain.WorkoutSet{},
		&domain.Metric{},
		&domain.DailySummary{},
		&domain.Goal{},
		&domain.Conversation{},
		&domain.Message{},
	); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	logger.Info("Database migrations completed")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationTime)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	router := httpAdapter.SetupRouter(authHandler, authService, cfg)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.String("address", fmt.Sprintf("http://localhost:%s", cfg.Server.Port)))

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
