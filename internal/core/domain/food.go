package domain

import (
	"time"

	"github.com/google/uuid"
)

// Food represents a food item with nutritional information
type Food struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FdcID       *int      `gorm:"uniqueIndex" json:"fdc_id,omitempty"` // USDA FoodData Central ID
	Name        string    `gorm:"type:varchar(255);not null;index" json:"name"` // Regular index for now, can add GIN full-text later
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Brand       *string   `gorm:"type:varchar(255)" json:"brand,omitempty"`
	Category    *string   `gorm:"type:varchar(100)" json:"category,omitempty"`

	// Base serving size (per 100g or 100ml)
	ServingSize     float64 `gorm:"type:decimal(10,2);not null" json:"serving_size"` // Stored as float64, precision 10,2
	ServingUnit     string  `gorm:"type:varchar(50);not null" json:"serving_unit"`

	// Macronutrients (per serving)
	Calories        float64  `gorm:"type:decimal(10,2);not null" json:"calories"`        // Stored as float64, precision 10,2
	Protein         float64  `gorm:"type:decimal(10,2);not null" json:"protein"`         // Stored as float64, precision 10,2
	Carbohydrates   float64  `gorm:"type:decimal(10,2);not null" json:"carbohydrates"`   // Stored as float64, precision 10,2
	Fat             float64  `gorm:"type:decimal(10,2);not null" json:"fat"`             // Stored as float64, precision 10,2
	Fiber           *float64 `gorm:"type:decimal(10,2)" json:"fiber,omitempty"`          // Stored as float64, precision 10,2
	Sugar           *float64 `gorm:"type:decimal(10,2)" json:"sugar,omitempty"`          // Stored as float64, precision 10,2
	SaturatedFat    *float64 `gorm:"type:decimal(10,2)" json:"saturated_fat,omitempty"`  // Stored as float64, precision 10,2
	TransFat        *float64 `gorm:"type:decimal(10,2)" json:"trans_fat,omitempty"`      // Stored as float64, precision 10,2
	Cholesterol     *float64 `gorm:"type:decimal(10,2)" json:"cholesterol,omitempty"`    // Stored as float64, precision 10,2
	Sodium          *float64 `gorm:"type:decimal(10,2)" json:"sodium,omitempty"`         // Stored as float64, precision 10,2
	Potassium       *float64 `gorm:"type:decimal(10,2)" json:"potassium,omitempty"`      // Stored as float64, precision 10,2

	// Metadata
	IsVerified bool       `gorm:"not null;default:false" json:"is_verified"`
	Source     *string    `gorm:"type:varchar(100)" json:"source,omitempty"` // e.g., "usda", "user", "manual"

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Ingredients          []FoodIngredient         `gorm:"foreignKey:FoodID" json:"ingredients,omitempty"`
	ServingConversions   []FoodServingConversion  `gorm:"foreignKey:FoodID" json:"serving_conversions,omitempty"`
}

// TableName specifies the table name for GORM
func (Food) TableName() string {
	return "foods"
}

// FoodIngredient represents an ingredient in a composite food
type FoodIngredient struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FoodID         uuid.UUID `gorm:"type:uuid;not null;index" json:"food_id"`
	IngredientID   uuid.UUID `gorm:"type:uuid;not null" json:"ingredient_id"`
	Quantity       float64   `gorm:"type:decimal(10,2);not null" json:"quantity"`      // Stored as float64, precision 10,2
	Unit           string    `gorm:"type:varchar(50);not null" json:"unit"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Food       Food `gorm:"foreignKey:FoodID" json:"-"`
	Ingredient Food `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

// TableName specifies the table name for GORM
func (FoodIngredient) TableName() string {
	return "food_ingredients"
}

// ServingUnit represents a common serving unit for conversions
type ServingUnit struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"` // e.g., "cup", "tablespoon", "piece"
	DisplayName string     `gorm:"type:varchar(100);not null" json:"display_name"`     // e.g., "Cup", "Tablespoon", "Piece"
	Category    string     `gorm:"type:varchar(50);not null" json:"category"`          // e.g., "volume", "weight", "count"

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for GORM
func (ServingUnit) TableName() string {
	return "serving_units"
}

// FoodServingConversion represents conversion factors for different serving sizes
type FoodServingConversion struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FoodID          uuid.UUID `gorm:"type:uuid;not null;index" json:"food_id"`
	ServingUnitID   uuid.UUID `gorm:"type:uuid;not null" json:"serving_unit_id"`
	GramsPerServing float64   `gorm:"type:decimal(10,2);not null" json:"grams_per_serving"` // Stored as float64, precision 10,2

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Food        Food        `gorm:"foreignKey:FoodID" json:"-"`
	ServingUnit ServingUnit `gorm:"foreignKey:ServingUnitID" json:"serving_unit,omitempty"`
}

// TableName specifies the table name for GORM
func (FoodServingConversion) TableName() string {
	return "food_serving_conversions"
}
