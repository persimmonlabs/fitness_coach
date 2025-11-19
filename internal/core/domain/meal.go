package domain

import (
	"time"

	"github.com/google/uuid"
)

// Meal represents a meal logged by a user
type Meal struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_meals" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	MealType  string    `gorm:"type:varchar(50);not null" json:"meal_type"` // breakfast, lunch, dinner, snack
	ConsumedAt time.Time `gorm:"not null;index:idx_user_meals" json:"consumed_at"`
	Notes     *string   `gorm:"type:text" json:"notes,omitempty"`

	// Calculated totals (denormalized for performance)
	TotalCalories     float64 `gorm:"type:decimal(10,2);not null" json:"total_calories"`      // Stored as float64, precision 10,2
	TotalProtein      float64 `gorm:"type:decimal(10,2);not null" json:"total_protein"`       // Stored as float64, precision 10,2
	TotalCarbohydrates float64 `gorm:"type:decimal(10,2);not null" json:"total_carbohydrates"` // Stored as float64, precision 10,2
	TotalFat          float64 `gorm:"type:decimal(10,2);not null" json:"total_fat"`           // Stored as float64, precision 10,2

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	FoodItems []MealFoodItem `gorm:"foreignKey:MealID" json:"food_items,omitempty"`
}

// TableName specifies the table name for GORM
func (Meal) TableName() string {
	return "meals"
}

// MealFoodItem represents a food item in a meal
type MealFoodItem struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MealID   uuid.UUID `gorm:"type:uuid;not null;index" json:"meal_id"`
	FoodID   uuid.UUID `gorm:"type:uuid;not null" json:"food_id"`
	Quantity float64   `gorm:"type:decimal(10,2);not null" json:"quantity"` // Stored as float64, precision 10,2
	Unit     string    `gorm:"type:varchar(50);not null" json:"unit"`

	// Calculated nutrition for this portion (denormalized)
	Calories      float64 `gorm:"type:decimal(10,2);not null" json:"calories"`      // Stored as float64, precision 10,2
	Protein       float64 `gorm:"type:decimal(10,2);not null" json:"protein"`       // Stored as float64, precision 10,2
	Carbohydrates float64 `gorm:"type:decimal(10,2);not null" json:"carbohydrates"` // Stored as float64, precision 10,2
	Fat           float64 `gorm:"type:decimal(10,2);not null" json:"fat"`           // Stored as float64, precision 10,2

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Meal Meal `gorm:"foreignKey:MealID" json:"-"`
	Food Food `gorm:"foreignKey:FoodID" json:"food,omitempty"`
}

// TableName specifies the table name for GORM
func (MealFoodItem) TableName() string {
	return "meal_food_items"
}

// NutritionTotals represents aggregated nutrition information
type NutritionTotals struct {
	TotalCalories      float64 `json:"total_calories"`
	TotalProtein       float64 `json:"total_protein"`
	TotalCarbohydrates float64 `json:"total_carbohydrates"`
	TotalFat           float64 `json:"total_fat"`
	TotalFiber         float64 `json:"total_fiber"`
	TotalSugar         float64 `json:"total_sugar"`
	TotalSodium        float64 `json:"total_sodium"`
}
