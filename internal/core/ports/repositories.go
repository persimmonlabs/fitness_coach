package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// FoodRepository defines the interface for food data operations
type FoodRepository interface {
	Create(ctx context.Context, food *domain.Food) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Food, error)
	GetByFdcID(ctx context.Context, fdcID int) (*domain.Food, error)
	Update(ctx context.Context, food *domain.Food) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Food, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.Food, error)

	// Ingredient operations
	AddIngredient(ctx context.Context, ingredient *domain.FoodIngredient) error
	GetIngredients(ctx context.Context, foodID uuid.UUID) ([]*domain.FoodIngredient, error)

	// Serving conversion operations
	AddServingConversion(ctx context.Context, conversion *domain.FoodServingConversion) error
	GetServingConversions(ctx context.Context, foodID uuid.UUID) ([]*domain.FoodServingConversion, error)

	// Serving unit operations
	CreateServingUnit(ctx context.Context, unit *domain.ServingUnit) error
	GetServingUnit(ctx context.Context, id uuid.UUID) (*domain.ServingUnit, error)
	ListServingUnits(ctx context.Context) ([]*domain.ServingUnit, error)
}

// MealRepository defines the interface for meal data operations
type MealRepository interface {
	Create(ctx context.Context, meal *domain.Meal) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Meal, error)
	Update(ctx context.Context, meal *domain.Meal) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*domain.Meal, error)

	// Food item operations
	AddFoodItem(ctx context.Context, item *domain.MealFoodItem) error
	UpdateFoodItem(ctx context.Context, item *domain.MealFoodItem) error
	RemoveFoodItem(ctx context.Context, id uuid.UUID) error
	GetFoodItems(ctx context.Context, mealID uuid.UUID) ([]*domain.MealFoodItem, error)
}

// ActivityRepository defines the interface for activity data operations
type ActivityRepository interface {
	Create(ctx context.Context, activity *domain.Activity) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Activity, error)
	Update(ctx context.Context, activity *domain.Activity) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*domain.Activity, error)
	GetTotalsByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
}

// WorkoutRepository defines the interface for workout data operations
type WorkoutRepository interface {
	Create(ctx context.Context, workout *domain.Workout) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Workout, error)
	Update(ctx context.Context, workout *domain.Workout) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*domain.Workout, error)

	// Exercise operations
	CreateExercise(ctx context.Context, exercise *domain.Exercise) error
	GetExercise(ctx context.Context, id uuid.UUID) (*domain.Exercise, error)
	ListExercises(ctx context.Context, category string, limit, offset int) ([]*domain.Exercise, error)

	// Workout exercise operations
	AddWorkoutExercise(ctx context.Context, workoutExercise *domain.WorkoutExercise) error
	GetWorkoutExercises(ctx context.Context, workoutID uuid.UUID) ([]*domain.WorkoutExercise, error)

	// Set operations
	AddSet(ctx context.Context, set *domain.WorkoutSet) error
	UpdateSet(ctx context.Context, set *domain.WorkoutSet) error
	DeleteSet(ctx context.Context, id uuid.UUID) error
	GetSets(ctx context.Context, workoutExerciseID uuid.UUID) ([]*domain.WorkoutSet, error)
}

// MetricRepository defines the interface for metric data operations
type MetricRepository interface {
	Create(ctx context.Context, metric *domain.Metric) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Metric, error)
	Update(ctx context.Context, metric *domain.Metric) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, metricType string, startDate, endDate time.Time, limit, offset int) ([]*domain.Metric, error)

	// Daily summary operations
	CreateOrUpdateDailySummary(ctx context.Context, summary *domain.DailySummary) error
	GetDailySummary(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.DailySummary, error)
	ListDailySummaries(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*domain.DailySummary, error)
}

// GoalRepository defines the interface for goal data operations
type GoalRepository interface {
	Create(ctx context.Context, goal *domain.Goal) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Goal, error)
	Update(ctx context.Context, goal *domain.Goal) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]*domain.Goal, error)
}

// ConversationRepository defines the interface for conversation data operations
type ConversationRepository interface {
	Create(ctx context.Context, conversation *domain.Conversation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Conversation, error)
	Update(ctx context.Context, conversation *domain.Conversation) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Conversation, error)

	// Message operations
	AddMessage(ctx context.Context, message *domain.Message) error
	GetMessages(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*domain.Message, error)
	GetLatestMessages(ctx context.Context, conversationID uuid.UUID, limit int) ([]*domain.Message, error)
}
