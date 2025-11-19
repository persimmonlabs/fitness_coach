package domain

import (
	"time"

	"github.com/google/uuid"
)

// ParsedMeal represents a meal that has been parsed from text or photo input
type ParsedMeal struct {
	MealType          string           `json:"meal_type"`
	LoggedAt          time.Time        `json:"logged_at"`
	FoodItems         []ParsedFoodItem `json:"food_items"`
	Confidence        float64          `json:"confidence"`
	NeedsConfirmation bool             `json:"needs_confirmation"`
}

// ParsedFoodItem represents a food item extracted from parsing
type ParsedFoodItem struct {
	FoodID      *uuid.UUID `json:"food_id,omitempty"`      // nil if AI-generated food
	FoodName    string     `json:"food_name"`
	Quantity    float64    `json:"quantity"`
	Unit        string     `json:"unit"`
	Confidence  float64    `json:"confidence"`
	AIGenerated bool       `json:"ai_generated"`
}

// NutritionEstimate represents AI-estimated nutrition information
type NutritionEstimate struct {
	CaloriesPer100g float64 `json:"calories_per_100g"`
	ProteinPer100g  float64 `json:"protein_per_100g"`
	CarbsPer100g    float64 `json:"carbs_per_100g"`
	FatPer100g      float64 `json:"fat_per_100g"`
	FiberPer100g    float64 `json:"fiber_per_100g"`
}
