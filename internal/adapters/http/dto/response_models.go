package dto

import "time"

// AuthResponse represents authentication response with user data and token
type AuthResponse struct {
	User  UserData `json:"user"`
	Token string   `json:"token"`
}

// UserData represents user information
type UserData struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// MealResponse represents a meal entry
type MealResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	MealType      string         `json:"meal_type"`
	ConsumedAt    time.Time      `json:"consumed_at"`
	Foods         []FoodItemResponse `json:"foods"`
	TotalCalories int            `json:"total_calories"`
	TotalProtein  float64        `json:"total_protein"`
	TotalCarbs    float64        `json:"total_carbs"`
	TotalFat      float64        `json:"total_fat"`
	Notes         string         `json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// FoodItemResponse represents a food item in a meal
type FoodItemResponse struct {
	FoodID   string  `json:"food_id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
}

// ParsedMealResponse represents AI-parsed meal data
type ParsedMealResponse struct {
	ParsedMealID  string         `json:"parsed_meal_id"`
	OriginalText  string         `json:"original_text"`
	MealType      string         `json:"meal_type"`
	DetectedFoods []DetectedFood `json:"detected_foods"`
	TotalCalories int            `json:"total_calories"`
	TotalProtein  float64        `json:"total_protein"`
	TotalCarbs    float64        `json:"total_carbs"`
	TotalFat      float64        `json:"total_fat"`
	Confidence    float64        `json:"confidence"`
}

// DetectedFood represents a food item detected by AI
type DetectedFood struct {
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
	Calories   int     `json:"calories"`
	Protein    float64 `json:"protein"`
	Carbs      float64 `json:"carbs"`
	Fat        float64 `json:"fat"`
	Confidence float64 `json:"confidence"`
}

// FoodResponse represents a food item
type FoodResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand,omitempty"`
	Calories    int       `json:"calories"`
	Protein     float64   `json:"protein"`
	Carbs       float64   `json:"carbs"`
	Fat         float64   `json:"fat"`
	Fiber       float64   `json:"fiber,omitempty"`
	Sugar       float64   `json:"sugar,omitempty"`
	ServingSize float64   `json:"serving_size"`
	ServingUnit string    `json:"serving_unit"`
	Barcode     string    `json:"barcode,omitempty"`
	IsCustom    bool      `json:"is_custom"`
	CreatedAt   time.Time `json:"created_at"`
}

// ActivityResponse represents an activity entry
type ActivityResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ActivityType string    `json:"activity_type"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     int       `json:"duration"` // in minutes
	Calories     int       `json:"calories"`
	Distance     float64   `json:"distance,omitempty"` // in km
	HeartRate    int       `json:"heart_rate,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	Source       string    `json:"source"` // manual, garmin, etc.
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// WorkoutResponse represents a workout session
type WorkoutResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	StartTime time.Time         `json:"start_time"`
	EndTime   *time.Time        `json:"end_time,omitempty"`
	Duration  int               `json:"duration,omitempty"` // in minutes
	Exercises []ExerciseSet     `json:"exercises"`
	Notes     string            `json:"notes,omitempty"`
	Status    string            `json:"status"` // in_progress, completed, cancelled
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ExerciseSet represents a set of exercises in a workout
type ExerciseSet struct {
	ExerciseID   string  `json:"exercise_id"`
	ExerciseName string  `json:"exercise_name"`
	SetNumber    int     `json:"set_number"`
	Reps         int     `json:"reps"`
	Weight       float64 `json:"weight,omitempty"`
	Duration     int     `json:"duration,omitempty"` // in seconds
	Notes        string  `json:"notes,omitempty"`
}

// MetricResponse represents a body metric entry
type MetricResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	MetricType string    `json:"metric_type"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	RecordedAt time.Time `json:"recorded_at"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// GoalResponse represents a fitness goal
type GoalResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	GoalType     string    `json:"goal_type"`
	TargetValue  float64   `json:"target_value"`
	CurrentValue float64   `json:"current_value"`
	Unit         string    `json:"unit"`
	Deadline     time.Time `json:"deadline"`
	Description  string    `json:"description,omitempty"`
	Status       string    `json:"status"` // active, completed, failed, cancelled
	Progress     float64   `json:"progress"` // percentage
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DailySummaryResponse represents a daily fitness summary
type DailySummaryResponse struct {
	Date             string  `json:"date"`
	TotalCalories    int     `json:"total_calories"`
	CalorieGoal      int     `json:"calorie_goal"`
	TotalProtein     float64 `json:"total_protein"`
	ProteinGoal      float64 `json:"protein_goal"`
	TotalCarbs       float64 `json:"total_carbs"`
	CarbsGoal        float64 `json:"carbs_goal"`
	TotalFat         float64 `json:"total_fat"`
	FatGoal          float64 `json:"fat_goal"`
	CaloriesBurned   int     `json:"calories_burned"`
	Activities       int     `json:"activities_count"`
	Workouts         int     `json:"workouts_count"`
	MealsLogged      int     `json:"meals_logged"`
	WaterIntake      float64 `json:"water_intake,omitempty"`
	StepsCount       int     `json:"steps_count,omitempty"`
	ActiveMinutes    int     `json:"active_minutes,omitempty"`
}

// ChatResponse represents AI coach response
type ChatResponse struct {
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
    Context   string    `json:"context,omitempty"`
    Suggestions []string `json:"suggestions,omitempty"`
}

// CompleteOnboardingResponse returns updated user and optional goal
type CompleteOnboardingResponse struct {
    User UserData      `json:"user"`
    Goal *GoalResponse `json:"goal,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// PaginatedResponse wraps paginated data
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ExerciseResponse represents an exercise
type ExerciseResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	MuscleGroup string   `json:"muscle_group"`
	Equipment   string   `json:"equipment,omitempty"`
	Description string   `json:"description,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
	IsCustom    bool     `json:"is_custom"`
	CreatedAt   time.Time `json:"created_at"`
}
