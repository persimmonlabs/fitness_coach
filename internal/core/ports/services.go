package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
)

// AuthService handles authentication and authorization
type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*domain.User, string, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
	GenerateJWT(userID string) (string, error)
	ParseJWT(token string) (string, error)
}

// FoodService handles food database operations
type FoodService interface {
	SearchFoods(ctx context.Context, query string, visibility *string, limit int) ([]*domain.Food, error)
	GetFood(ctx context.Context, foodID string) (*domain.Food, error)
	CreateFood(ctx context.Context, food *domain.Food) (*domain.Food, error)
	UpdateFood(ctx context.Context, foodID string, updates map[string]interface{}) (*domain.Food, error)
	CreateAIGeneratedFood(ctx context.Context, name string, nutritionData map[string]interface{}) (*domain.Food, error)
}

// MealService handles meal tracking and nutrition calculation
type MealService interface {
	GetMeals(ctx context.Context, userID string, date *time.Time) ([]*domain.Meal, error)
	GetMeal(ctx context.Context, mealID string) (*domain.Meal, error)
	CreateMeal(ctx context.Context, userID string, mealData *domain.Meal) (*domain.Meal, error)
	ConfirmParsedMeal(ctx context.Context, userID string, parsedMeal *domain.Meal) (*domain.Meal, error)
	UpdateMeal(ctx context.Context, mealID string, updates map[string]interface{}) (*domain.Meal, error)
	DeleteMeal(ctx context.Context, mealID string) error
	CalculateMealNutrition(ctx context.Context, mealID string) (*domain.NutritionTotals, error)
}

// ActivityService handles activity tracking
type ActivityService interface {
	GetActivities(ctx context.Context, userID string, startDate, endDate *time.Time) ([]*domain.Activity, error)
	GetActivity(ctx context.Context, activityID string) (*domain.Activity, error)
	CreateActivity(ctx context.Context, userID string, activityData *domain.Activity) (*domain.Activity, error)
	UpdateActivity(ctx context.Context, activityID string, updates map[string]interface{}) (*domain.Activity, error)
	DeleteActivity(ctx context.Context, activityID string) error
}

// WorkoutService handles workout tracking
type WorkoutService interface {
	StartWorkout(ctx context.Context, userID, name string) (*domain.Workout, error)
	GetWorkouts(ctx context.Context, userID string, startDate, endDate *time.Time) ([]*domain.Workout, error)
	GetWorkout(ctx context.Context, workoutID string) (*domain.Workout, error)
	AddExercise(ctx context.Context, workoutID, exerciseID string) (*domain.WorkoutExercise, error)
	LogSet(ctx context.Context, workoutExerciseID string, setData *domain.WorkoutSet) (*domain.WorkoutSet, error)
	FinishWorkout(ctx context.Context, workoutID string) error
	DeleteWorkout(ctx context.Context, workoutID string) error
}

// MetricService handles health metrics tracking
type MetricService interface {
	LogMetric(ctx context.Context, userID, metricType string, value float64, unit string, recordedAt time.Time) (*domain.Metric, error)
	GetMetricTrend(ctx context.Context, userID, metricType string, startDate, endDate *time.Time) ([]*domain.Metric, error)
	GetLatestMetric(ctx context.Context, userID, metricType string) (*domain.Metric, error)
}

// GoalService handles user goals
type GoalService interface {
	CreateGoal(ctx context.Context, userID string, goalData *domain.Goal) (*domain.Goal, error)
	GetGoals(ctx context.Context, userID string, status *string) ([]*domain.Goal, error)
	UpdateGoal(ctx context.Context, goalID string, updates map[string]interface{}) (*domain.Goal, error)
	DeleteGoal(ctx context.Context, goalID string) error
}

// SummaryService handles daily summaries
type SummaryService interface {
	GetDailySummary(ctx context.Context, userID string, date time.Time) (*domain.DailySummary, error)
	CalculateDailySummary(ctx context.Context, userID string, date time.Time) (*domain.DailySummary, error)
}

// AgentResponse represents the response from the AI agent
type AgentResponse struct {
	Message    string    `json:"message"`
	ToolsUsed  []string  `json:"tools_used"`
	Confidence float64   `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
}

// AgentService handles AI agent interactions
type AgentService interface {
	SendMessage(ctx context.Context, userID uuid.UUID, message string) (*AgentResponse, error)
}
