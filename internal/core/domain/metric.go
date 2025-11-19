package domain

import (
	"time"

	"github.com/google/uuid"
)

// Metric represents a health/fitness metric measurement
type Metric struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index:idx_user_metrics" json:"user_id"`
	MetricType string    `gorm:"type:varchar(100);not null;index:idx_user_metrics" json:"metric_type"` // weight, body_fat, steps, etc.
	Value      float64   `gorm:"type:decimal(10,2);not null" json:"value"` // Stored as float64, precision 10,2
	Unit       string    `gorm:"type:varchar(50);not null" json:"unit"`
	MeasuredAt time.Time `gorm:"not null;index:idx_user_metrics" json:"measured_at"`
	Notes      *string   `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for GORM
func (Metric) TableName() string {
	return "metrics"
}

// DailySummary represents aggregated daily health data
type DailySummary struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_date" json:"user_id"`
	Date   time.Time `gorm:"type:date;not null;uniqueIndex:idx_user_date" json:"date"`

	// Nutrition totals
	TotalCalories     float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total_calories"`      // Stored as float64, precision 10,2
	TotalProtein      float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total_protein"`       // Stored as float64, precision 10,2
	TotalCarbohydrates float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total_carbohydrates"` // Stored as float64, precision 10,2
	TotalFat          float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total_fat"`           // Stored as float64, precision 10,2

	// Activity totals
	TotalCaloriesBurned  float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total_calories_burned"` // Stored as float64, precision 10,2
	TotalExerciseMinutes int     `gorm:"type:integer;not null;default:0" json:"total_exercise_minutes"`
	TotalSteps           int     `gorm:"type:integer;not null;default:0" json:"total_steps"`
	TotalDistance        float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total_distance"` // Stored as float64, precision 10,2, in km

	// Metrics
	Weight   *float64 `gorm:"type:decimal(5,2)" json:"weight,omitempty"`    // Stored as float64, precision 5,2, in kg
	BodyFat  *float64 `gorm:"type:decimal(5,2)" json:"body_fat,omitempty"`  // Stored as float64, precision 5,2, percentage

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for GORM
func (DailySummary) TableName() string {
	return "daily_summaries"
}
