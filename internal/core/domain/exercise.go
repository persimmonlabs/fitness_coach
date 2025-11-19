package domain

import (
	"time"

	"github.com/google/uuid"
)

// Exercise represents a type of exercise (e.g., bench press, squat)
type Exercise struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Category    string    `gorm:"type:varchar(100);not null" json:"category"` // strength, cardio, flexibility, etc.
	MuscleGroup *string   `gorm:"type:varchar(100)" json:"muscle_group,omitempty"` // chest, legs, back, etc.
	Equipment   *string   `gorm:"type:varchar(100)" json:"equipment,omitempty"` // barbell, dumbbell, bodyweight, etc.
	Difficulty  *string   `gorm:"type:varchar(50)" json:"difficulty,omitempty"` // beginner, intermediate, advanced

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for GORM
func (Exercise) TableName() string {
	return "exercises"
}
