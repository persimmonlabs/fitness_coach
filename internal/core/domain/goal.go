package domain

import (
	"time"

	"github.com/google/uuid"
)

// Goal represents a user's fitness or health goal
type Goal struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_user_goals" json:"user_id"`
	GoalType    string     `gorm:"type:varchar(100);not null" json:"goal_type"` // weight_loss, muscle_gain, distance, etc.
	Description string     `gorm:"type:text;not null" json:"description"`

	TargetValue   float64    `gorm:"type:decimal(10,2);not null" json:"target_value"` // Stored as float64, precision 10,2
	CurrentValue  *float64   `gorm:"type:decimal(10,2)" json:"current_value,omitempty"` // Stored as float64, precision 10,2
	Unit          string     `gorm:"type:varchar(50);not null" json:"unit"`

	StartDate     time.Time  `gorm:"not null" json:"start_date"`
	TargetDate    *time.Time `json:"target_date,omitempty"`
	CompletedDate *time.Time `json:"completed_date,omitempty"`

	Status        string     `gorm:"type:varchar(50);not null;default:'active'" json:"status"` // active, completed, abandoned

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for GORM
func (Goal) TableName() string {
	return "goals"
}
