package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"` // Never expose password hash in JSON
	FirstName    string     `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName     string     `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	DateOfBirth  *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Gender       *string    `gorm:"type:varchar(20)" json:"gender,omitempty"`
	HeightCm     *float64   `gorm:"type:decimal(5,2)" json:"height_cm,omitempty"` // Stored as float64, precision documented
	WeightKg     *float64   `gorm:"type:decimal(5,2)" json:"weight_kg,omitempty"` // Stored as float64, precision documented
	ActivityLevel *string   `gorm:"type:varchar(50)" json:"activity_level,omitempty"`

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
