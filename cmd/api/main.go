package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	httpAdapter "fitness-tracker/internal/adapters/http"
	"fitness-tracker/internal/adapters/http/handlers"
	"fitness-tracker/internal/adapters/llm"
	"fitness-tracker/internal/adapters/persistence"
	"fitness-tracker/internal/core/services"
	"fitness-tracker/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Initialize logger
	logger, err := initLogger(cfg.Server.Environment)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Fitness Coach API",
		zap.String("environment", cfg.Server.Environment),
		zap.String("version", "1.0.0"),
	)

	// Connect to database
	db, err := connectDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	if err := runMigrations(db, logger); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize dependencies
	deps := initializeDependencies(db, cfg, logger)

	// Setup router
	router := httpAdapter.SetupRouter(
		deps.authHandler,
		deps.mealHandler,
		deps.foodHandler,
		deps.activityHandler,
		deps.workoutHandler,
		deps.exerciseHandler,
		deps.metricHandler,
		deps.goalHandler,
		deps.chatHandler,
		deps.summaryHandler,
		cfg,
		logger,
	)

	// Start HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// initLogger initializes the zap logger based on environment
func initLogger(environment string) (*zap.Logger, error) {
	if environment == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

// connectDatabase establishes database connection
func connectDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connection established",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Name),
	)

	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB, logger *zap.Logger) error {
	logger.Info("Running database migrations...")

	// Auto migrate all models
	err := db.AutoMigrate(
		&persistence.UserModel{},
		&persistence.MealModel{},
		&persistence.FoodModel{},
		&persistence.ActivityModel{},
		&persistence.WorkoutModel{},
		&persistence.ExerciseModel{},
		&persistence.MetricModel{},
		&persistence.GoalModel{},
		&persistence.ChatMessageModel{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// dependencies holds all application dependencies
type dependencies struct {
	authHandler     *handlers.AuthHandler
	mealHandler     *handlers.MealHandler
	foodHandler     *handlers.FoodHandler
	activityHandler *handlers.ActivityHandler
	workoutHandler  *handlers.WorkoutHandler
	exerciseHandler *handlers.ExerciseHandler
	metricHandler   *handlers.MetricHandler
	goalHandler     *handlers.GoalHandler
	chatHandler     *handlers.ChatHandler
	summaryHandler  *handlers.SummaryHandler
}

// initializeDependencies initializes all application dependencies using dependency injection
func initializeDependencies(db *gorm.DB, cfg *config.Config, logger *zap.Logger) *dependencies {
	// Initialize repositories
	userRepo := persistence.NewUserRepository(db)
	mealRepo := persistence.NewMealRepository(db)
	foodRepo := persistence.NewFoodRepository(db)
	activityRepo := persistence.NewActivityRepository(db)
	workoutRepo := persistence.NewWorkoutRepository(db)
	exerciseRepo := persistence.NewExerciseRepository(db)
	metricRepo := persistence.NewMetricRepository(db)
	goalRepo := persistence.NewGoalRepository(db)
	chatRepo := persistence.NewChatRepository(db)

	// Initialize LLM client
	llmClient := llm.NewOpenAIClient(cfg.LLM.APIKey, cfg.LLM.Model, logger)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	mealService := services.NewMealService(mealRepo)
	foodService := services.NewFoodService(foodRepo)
	activityService := services.NewActivityService(activityRepo)
	workoutService := services.NewWorkoutService(workoutRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	metricService := services.NewMetricService(metricRepo)
	goalService := services.NewGoalService(goalRepo)
	chatService := services.NewChatService(chatRepo, llmClient)
	summaryService := services.NewSummaryService(
		mealRepo,
		activityRepo,
		workoutRepo,
		metricRepo,
		goalRepo,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	mealHandler := handlers.NewMealHandler(mealService, logger)
	foodHandler := handlers.NewFoodHandler(foodService, logger)
	activityHandler := handlers.NewActivityHandler(activityService, logger)
	workoutHandler := handlers.NewWorkoutHandler(workoutService, logger)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService, logger)
	metricHandler := handlers.NewMetricHandler(metricService, logger)
	goalHandler := handlers.NewGoalHandler(goalService, logger)
	chatHandler := handlers.NewChatHandler(chatService, logger)
	summaryHandler := handlers.NewSummaryHandler(summaryService, logger)

	logger.Info("All dependencies initialized successfully")

	return &dependencies{
		authHandler:     authHandler,
		mealHandler:     mealHandler,
		foodHandler:     foodHandler,
		activityHandler: activityHandler,
		workoutHandler:  workoutHandler,
		exerciseHandler: exerciseHandler,
		metricHandler:   metricHandler,
		goalHandler:     goalHandler,
		chatHandler:     chatHandler,
		summaryHandler:  summaryHandler,
	}
}
