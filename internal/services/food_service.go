package services

import (
	"context"
	"fmt"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"

	"github.com/google/uuid"
)

type foodService struct {
	foodRepo ports.FoodRepository
}

// NewFoodService creates a new food service
func NewFoodService(foodRepo ports.FoodRepository) ports.FoodService {
	return &foodService{
		foodRepo: foodRepo,
	}
}

func (s *foodService) SearchFoods(ctx context.Context, query string, visibility *string, limit int) ([]*domain.Food, error) {
	if limit <= 0 {
		limit = 20 // default limit
	}
	if limit > 100 {
		limit = 100 // max limit
	}

	foods, err := s.foodRepo.Search(ctx, query, visibility, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search foods: %w", err)
	}

	return foods, nil
}

func (s *foodService) GetFood(ctx context.Context, foodID string) (*domain.Food, error) {
	if foodID == "" {
		return nil, domain.ErrInvalidInput
	}

	food, err := s.foodRepo.GetByID(ctx, foodID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get food: %w", err)
	}

	return food, nil
}

func (s *foodService) CreateFood(ctx context.Context, food *domain.Food) (*domain.Food, error) {
	// Validate food data
	if food == nil {
		return nil, domain.ErrInvalidInput
	}
	if food.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Set ID if not provided
	if food.ID == "" {
		food.ID = uuid.New().String()
	}

	// Default visibility to private
	if food.Visibility == "" {
		food.Visibility = "private"
	}

	// Validate visibility
	if food.Visibility != "private" && food.Visibility != "public" {
		return nil, domain.ErrInvalidInput
	}

	// Create food
	if err := s.foodRepo.Create(ctx, food); err != nil {
		return nil, fmt.Errorf("failed to create food: %w", err)
	}

	return food, nil
}

func (s *foodService) UpdateFood(ctx context.Context, foodID string, updates map[string]interface{}) (*domain.Food, error) {
	if foodID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify food exists
	existing, err := s.foodRepo.GetByID(ctx, foodID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get food: %w", err)
	}

	// Validate visibility if being updated
	if visibility, ok := updates["visibility"].(string); ok {
		if visibility != "private" && visibility != "public" {
			return nil, domain.ErrInvalidInput
		}
	}

	// Update food
	if err := s.foodRepo.Update(ctx, foodID, updates); err != nil {
		return nil, fmt.Errorf("failed to update food: %w", err)
	}

	// Return updated food
	return s.foodRepo.GetByID(ctx, existing.ID)
}

func (s *foodService) CreateAIGeneratedFood(ctx context.Context, name string, nutritionData map[string]interface{}) (*domain.Food, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	food := &domain.Food{
		ID:         uuid.New().String(),
		Name:       name,
		Source:     "ai_generated",
		Visibility: "private",
	}

	// Extract nutrition data
	if calories, ok := nutritionData["calories"].(float64); ok {
		food.Calories = calories
	}
	if protein, ok := nutritionData["protein"].(float64); ok {
		food.Protein = protein
	}
	if carbs, ok := nutritionData["carbs"].(float64); ok {
		food.Carbs = carbs
	}
	if fat, ok := nutritionData["fat"].(float64); ok {
		food.Fat = fat
	}
	if fiber, ok := nutritionData["fiber"].(float64); ok {
		food.Fiber = fiber
	}
	if sugar, ok := nutritionData["sugar"].(float64); ok {
		food.Sugar = sugar
	}
	if sodium, ok := nutritionData["sodium"].(float64); ok {
		food.Sodium = sodium
	}

	// Create food
	if err := s.foodRepo.Create(ctx, food); err != nil {
		return nil, fmt.Errorf("failed to create AI generated food: %w", err)
	}

	return food, nil
}
