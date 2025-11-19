package domain

import (
	"time"

	"github.com/google/uuid"
)

// Workout represents a workout session
type Workout struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_workouts" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	StartTime time.Time `gorm:"not null;index:idx_user_workouts" json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`

	DurationMinutes    *int     `gorm:"type:integer" json:"duration_minutes,omitempty"`
	CaloriesBurned     *float64 `gorm:"type:decimal(10,2)" json:"calories_burned,omitempty"` // Stored as float64, precision 10,2
	AverageHeartRate   *int     `gorm:"type:integer" json:"average_heart_rate,omitempty"`
	MaxHeartRate       *int     `gorm:"type:integer" json:"max_heart_rate,omitempty"`

	Notes *string `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User      User              `gorm:"foreignKey:UserID" json:"-"`
	Exercises []WorkoutExercise `gorm:"foreignKey:WorkoutID" json:"exercises,omitempty"`
}

// TableName specifies the table name for GORM
func (Workout) TableName() string {
	return "workouts"
}

// WorkoutExercise represents an exercise performed in a workout
type WorkoutExercise struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkoutID  uuid.UUID `gorm:"type:uuid;not null;index" json:"workout_id"`
	ExerciseID uuid.UUID `gorm:"type:uuid;not null" json:"exercise_id"`
	OrderIndex int       `gorm:"type:integer;not null" json:"order_index"` // Order in the workout
	Notes      *string   `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Workout  Workout      `gorm:"foreignKey:WorkoutID" json:"-"`
	Exercise Exercise     `gorm:"foreignKey:ExerciseID" json:"exercise,omitempty"`
	Sets     []WorkoutSet `gorm:"foreignKey:WorkoutExerciseID" json:"sets,omitempty"`
}

// TableName specifies the table name for GORM
func (WorkoutExercise) TableName() string {
	return "workout_exercises"
}

// WorkoutSet represents a set in an exercise
type WorkoutSet struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkoutExerciseID  uuid.UUID `gorm:"type:uuid;not null;index" json:"workout_exercise_id"`
	SetNumber          int       `gorm:"type:integer;not null" json:"set_number"`

	Reps               *int     `gorm:"type:integer" json:"reps,omitempty"`
	Weight             *float64 `gorm:"type:decimal(10,2)" json:"weight,omitempty"`      // Stored as float64, precision 10,2, in kg
	DurationSeconds    *int     `gorm:"type:integer" json:"duration_seconds,omitempty"`
	Distance           *float64 `gorm:"type:decimal(10,2)" json:"distance,omitempty"`    // Stored as float64, precision 10,2, in meters
	RestSeconds        *int     `gorm:"type:integer" json:"rest_seconds,omitempty"`

	Notes *string `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	WorkoutExercise WorkoutExercise `gorm:"foreignKey:WorkoutExerciseID" json:"-"`
}

// TableName specifies the table name for GORM
func (WorkoutSet) TableName() string {
	return "workout_sets"
}
