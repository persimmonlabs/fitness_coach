package dto

import "time"

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

// LoginRequest represents user login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CreateMealRequest represents a new meal entry
type CreateMealRequest struct {
	Name        string    `json:"name" validate:"required"`
	MealType    string    `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	ConsumedAt  time.Time `json:"consumed_at" validate:"required"`
	Foods       []FoodItem `json:"foods" validate:"required,dive"`
	TotalCalories int     `json:"total_calories,omitempty"`
	TotalProtein  float64 `json:"total_protein,omitempty"`
	TotalCarbs    float64 `json:"total_carbs,omitempty"`
	TotalFat      float64 `json:"total_fat,omitempty"`
	Notes       string    `json:"notes,omitempty"`
}

// FoodItem represents a food item in a meal
type FoodItem struct {
	FoodID   string  `json:"food_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
	Unit     string  `json:"unit" validate:"required"`
}

// ParseMealRequest represents natural language meal input
type ParseMealRequest struct {
	Description string    `json:"description" validate:"required"`
	MealType    string    `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	ConsumedAt  time.Time `json:"consumed_at,omitempty"`
}

// ConfirmMealRequest confirms a parsed meal
type ConfirmMealRequest struct {
	ParsedMealID string `json:"parsed_meal_id" validate:"required"`
	Adjustments  *CreateMealRequest `json:"adjustments,omitempty"`
}

// CreateFoodRequest represents a new custom food entry
type CreateFoodRequest struct {
	Name        string  `json:"name" validate:"required"`
	Brand       string  `json:"brand,omitempty"`
	Calories    int     `json:"calories" validate:"required,gte=0"`
	Protein     float64 `json:"protein" validate:"required,gte=0"`
	Carbs       float64 `json:"carbs" validate:"required,gte=0"`
	Fat         float64 `json:"fat" validate:"required,gte=0"`
	Fiber       float64 `json:"fiber,omitempty"`
	Sugar       float64 `json:"sugar,omitempty"`
	ServingSize float64 `json:"serving_size" validate:"required,gt=0"`
	ServingUnit string  `json:"serving_unit" validate:"required"`
	Barcode     string  `json:"barcode,omitempty"`
}

// CreateActivityRequest represents a new activity entry
type CreateActivityRequest struct {
	ActivityType string    `json:"activity_type" validate:"required"`
	StartTime    time.Time `json:"start_time" validate:"required"`
	EndTime      time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Duration     int       `json:"duration,omitempty"` // in minutes
	Calories     int       `json:"calories,omitempty"`
	Distance     float64   `json:"distance,omitempty"` // in km
	HeartRate    int       `json:"heart_rate,omitempty"`
	Notes        string    `json:"notes,omitempty"`
}

// StartWorkoutRequest represents starting a new workout
type StartWorkoutRequest struct {
	Name      string    `json:"name" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	StartTime time.Time `json:"start_time,omitempty"`
	Notes     string    `json:"notes,omitempty"`
}

// LogSetRequest represents logging a set during a workout
type LogSetRequest struct {
	WorkoutID  string  `json:"workout_id" validate:"required"`
	ExerciseID string  `json:"exercise_id" validate:"required"`
	SetNumber  int     `json:"set_number" validate:"required,gt=0"`
	Reps       int     `json:"reps" validate:"required,gt=0"`
	Weight     float64 `json:"weight,omitempty"`
	Duration   int     `json:"duration,omitempty"` // in seconds
	Notes      string  `json:"notes,omitempty"`
}

// LogMetricRequest represents logging a body metric
type LogMetricRequest struct {
	MetricType string    `json:"metric_type" validate:"required,oneof=weight body_fat muscle_mass bmi waist_circumference"`
	Value      float64   `json:"value" validate:"required,gt=0"`
	Unit       string    `json:"unit" validate:"required"`
	RecordedAt time.Time `json:"recorded_at,omitempty"`
	Notes      string    `json:"notes,omitempty"`
}

// CreateGoalRequest represents a new fitness goal
type CreateGoalRequest struct {
	GoalType    string    `json:"goal_type" validate:"required,oneof=weight calorie_intake protein_intake workout_frequency custom"`
	TargetValue float64   `json:"target_value" validate:"required"`
	CurrentValue float64  `json:"current_value,omitempty"`
	Unit        string    `json:"unit" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	Description string    `json:"description,omitempty"`
}

// ChatRequest represents a chat message to the AI coach
type ChatRequest struct {
    Message string `json:"message" validate:"required,min=1"`
    Context string `json:"context,omitempty"` // Additional context for the AI
}

// CompleteOnboardingRequest captures onboarding profile fields
type CompleteOnboardingRequest struct {
    Age                int                    `json:"age" validate:"required,min=13,max=120"`
    Sex                string                 `json:"sex" validate:"required,oneof=male female other"`
    HeightCm           float64                `json:"height_cm" validate:"required,gt=0"`
    CurrentWeight      float64                `json:"current_weight" validate:"required,gt=0"`
    ActivityLevel      string                 `json:"activity_level" validate:"required"`
    GoalType           string                 `json:"goal_type,omitempty"`
    TargetWeight       float64                `json:"target_weight,omitempty"`
    TargetDate         time.Time              `json:"target_date,omitempty"`
    UnitSystem         string                 `json:"unit_system" validate:"required,oneof=metric imperial"`
    Timezone           string                 `json:"timezone" validate:"required"`
    DietaryPreferences map[string]interface{} `json:"dietary_preferences,omitempty"`
}
