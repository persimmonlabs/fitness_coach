package domain

import (
	"time"

	"github.com/google/uuid"
)

// Activity represents a physical activity logged by a user
type Activity struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index:idx_user_activities" json:"user_id"`
	ActivityType  string    `gorm:"type:varchar(100);not null" json:"activity_type"` // walking, running, cycling, etc.
	StartTime     time.Time `gorm:"not null;index:idx_user_activities" json:"start_time"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	DurationMinutes *int    `gorm:"type:integer" json:"duration_minutes,omitempty"`

	// Activity metrics
	Distance       *float64 `gorm:"type:decimal(10,2)" json:"distance,omitempty"`        // Stored as float64, precision 10,2, in km
	CaloriesBurned *float64 `gorm:"type:decimal(10,2)" json:"calories_burned,omitempty"` // Stored as float64, precision 10,2
	AverageHeartRate *int   `gorm:"type:integer" json:"average_heart_rate,omitempty"`
	MaxHeartRate   *int     `gorm:"type:integer" json:"max_heart_rate,omitempty"`
	Steps          *int     `gorm:"type:integer" json:"steps,omitempty"`

	Notes *string `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for GORM
func (Activity) TableName() string {
	return "activities"
}
