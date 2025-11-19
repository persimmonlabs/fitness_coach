package integration

import (
	"context"
	"testing"
	"time"

	"fitness-tracker/internal/core/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDatabase holds the test database container and connection
type TestDatabase struct {
	Container *postgres.PostgresContainer
	DB        *gorm.DB
}

// SetupTestDB creates a PostgreSQL container, runs migrations, and returns a GORM DB instance
func SetupTestDB(t *testing.T) *TestDatabase {
	ctx := context.Background()

	// Create PostgreSQL container
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_fitness_tracker"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get connection string")

	// Connect to database
	db, err := gorm.Open(pgdriver.Open(connStr), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Run migrations
	err = runMigrations(db)
	require.NoError(t, err, "Failed to run migrations")

	return &TestDatabase{
		Container: pgContainer,
		DB:        db,
	}
}

// TeardownTestDB cleans up the test database container
func TeardownTestDB(t *testing.T, testDB *TestDatabase) {
	if testDB == nil {
		return
	}

	ctx := context.Background()

	// Close database connection
	if testDB.DB != nil {
		sqlDB, err := testDB.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	// Terminate container
	if testDB.Container != nil {
		err := testDB.Container.Terminate(ctx)
		require.NoError(t, err, "Failed to terminate container")
	}
}

// runMigrations applies all database migrations
func runMigrations(db *gorm.DB) error {
	// Auto-migrate all domain models
	return db.AutoMigrate(
		&domain.User{},
		&domain.Food{},
		&domain.ServingUnit{},
		&domain.FoodIngredient{},
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
	)
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *gorm.DB, email string) *domain.User {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("test_password123"), bcrypt.DefaultCost)
	require.NoError(t, err, "Failed to hash password")

	user := &domain.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		FirstName:    "Test",
		LastName:     "User",
	}

	err = db.Create(user).Error
	require.NoError(t, err, "Failed to create test user")

	return user
}

// GetTestJWT generates a JWT token for a test user
func GetTestJWT(userID uuid.UUID) (string, error) {
	// Use a test secret key
	secretKey := []byte("test_secret_key_for_integration_tests")

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// CreateTestFood creates a test food item in the database
func CreateTestFood(t *testing.T, db *gorm.DB, name string, calories float64) *domain.Food {
	food := &domain.Food{
		Name:          name,
		ServingSize:   100,
		ServingUnit:   "g",
		Calories:      calories,
		Protein:       10.0,
		Carbohydrates: 20.0,
		Fat:           5.0,
		IsVerified:    true,
	}

	err := db.Create(food).Error
	require.NoError(t, err, "Failed to create test food")

	return food
}

// CreateTestExercise creates a test exercise in the database
func CreateTestExercise(t *testing.T, db *gorm.DB, name, category string) *domain.Exercise {
	exercise := &domain.Exercise{
		Name:        name,
		Category:    category,
		Description: stringPtr("Test exercise"),
		MuscleGroup: stringPtr("Test muscle"),
	}

	err := db.Create(exercise).Error
	require.NoError(t, err, "Failed to create test exercise")

	return exercise
}

// CreateTestMeal creates a test meal in the database
func CreateTestMeal(t *testing.T, db *gorm.DB, userID uuid.UUID, mealType string) *domain.Meal {
	meal := &domain.Meal{
		UserID:             userID,
		Name:               "Test Meal",
		MealType:           mealType,
		ConsumedAt:         time.Now(),
		TotalCalories:      500,
		TotalProtein:       30,
		TotalCarbohydrates: 60,
		TotalFat:           15,
	}

	err := db.Create(meal).Error
	require.NoError(t, err, "Failed to create test meal")

	return meal
}

// CreateTestActivity creates a test activity in the database
func CreateTestActivity(t *testing.T, db *gorm.DB, userID uuid.UUID, activityType string) *domain.Activity {
	duration := 30
	calories := 200.0

	activity := &domain.Activity{
		UserID:          userID,
		ActivityType:    activityType,
		StartTime:       time.Now().Add(-1 * time.Hour),
		DurationMinutes: &duration,
		CaloriesBurned:  &calories,
	}

	err := db.Create(activity).Error
	require.NoError(t, err, "Failed to create test activity")

	return activity
}

// CleanupTestData removes all test data from the database
func CleanupTestData(t *testing.T, db *gorm.DB) {
	// Delete in reverse order of dependencies
	db.Exec("TRUNCATE TABLE messages CASCADE")
	db.Exec("TRUNCATE TABLE conversations CASCADE")
	db.Exec("TRUNCATE TABLE goals CASCADE")
	db.Exec("TRUNCATE TABLE daily_summaries CASCADE")
	db.Exec("TRUNCATE TABLE metrics CASCADE")
	db.Exec("TRUNCATE TABLE workout_sets CASCADE")
	db.Exec("TRUNCATE TABLE workout_exercises CASCADE")
	db.Exec("TRUNCATE TABLE workouts CASCADE")
	db.Exec("TRUNCATE TABLE exercises CASCADE")
	db.Exec("TRUNCATE TABLE activities CASCADE")
	db.Exec("TRUNCATE TABLE meal_food_items CASCADE")
	db.Exec("TRUNCATE TABLE meals CASCADE")
	db.Exec("TRUNCATE TABLE food_serving_conversions CASCADE")
	db.Exec("TRUNCATE TABLE food_ingredients CASCADE")
	db.Exec("TRUNCATE TABLE serving_units CASCADE")
	db.Exec("TRUNCATE TABLE foods CASCADE")
	db.Exec("TRUNCATE TABLE users CASCADE")
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// int64Ptr returns a pointer to an int64
func int64Ptr(i int64) *int64 {
	return &i
}

// float64Ptr returns a pointer to a float64
func float64Ptr(f float64) *float64 {
	return &f
}

// intPtr returns a pointer to an int
func intPtr(i int) *int {
	return &i
}
