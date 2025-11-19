package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection and runs migrations
func InitDB(config *DatabaseConfig) (*gorm.DB, error) {
	// Configure GORM logger
	gormLogger := logger.Default
	if config.SSLMode == "disable" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database at %s:%d", config.Host, config.Port)

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")

	return db, nil
}

// runMigrations runs automatic database migrations
func runMigrations(db *gorm.DB) error {
	// Import models here when they are created
	// For now, we'll define the migration interface

	// Example migration structure (uncomment when models are ready):
	/*
	err := db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Goal{},
		&models.Workout{},
		&models.Exercise{},
		&models.WorkoutExercise{},
		&models.WorkoutLog{},
		&models.ExerciseLog{},
		&models.Meal{},
		&models.MealFood{},
		&models.Food{},
		&models.NutritionLog{},
		&models.ProgressMetric{},
		&models.CoachConversation{},
		&models.CoachMessage{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}
	*/

	// Create extensions if needed
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: failed to create uuid-ossp extension: %v", err)
	}

	// Create custom types if needed
	createCustomTypes(db)

	return nil
}

// createCustomTypes creates custom PostgreSQL types
func createCustomTypes(db *gorm.DB) {
	// Create ENUM types for various fields
	customTypes := []string{
		`DO $$ BEGIN
			CREATE TYPE user_role AS ENUM ('user', 'admin', 'coach');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		`DO $$ BEGIN
			CREATE TYPE goal_type AS ENUM ('weight_loss', 'muscle_gain', 'maintenance', 'endurance', 'flexibility', 'general_fitness');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		`DO $$ BEGIN
			CREATE TYPE goal_status AS ENUM ('active', 'completed', 'abandoned', 'paused');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		`DO $$ BEGIN
			CREATE TYPE activity_level AS ENUM ('sedentary', 'lightly_active', 'moderately_active', 'very_active', 'extremely_active');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		`DO $$ BEGIN
			CREATE TYPE gender AS ENUM ('male', 'female', 'other', 'prefer_not_to_say');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		`DO $$ BEGIN
			CREATE TYPE exercise_type AS ENUM ('strength', 'cardio', 'flexibility', 'balance', 'sports', 'other');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		`DO $$ BEGIN
			CREATE TYPE difficulty_level AS ENUM ('beginner', 'intermediate', 'advanced', 'expert');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		`DO $$ BEGIN
			CREATE TYPE message_sender AS ENUM ('user', 'coach');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,
	}

	for _, typeSQL := range customTypes {
		if err := db.Exec(typeSQL).Error; err != nil {
			log.Printf("Warning: failed to create custom type: %v", err)
		}
	}
}

// CloseDB closes the database connection gracefully
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed successfully")
	return nil
}

// HealthCheck performs a database health check
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetDBStats returns database connection statistics
func GetDBStats(db *gorm.DB) (map[string]interface{}, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}
