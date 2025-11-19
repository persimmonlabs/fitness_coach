package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type mealService struct {
	mealRepo ports.MealRepository
	foodRepo ports.FoodRepository
}

// NewMealService creates a new meal service
func NewMealService(mealRepo ports.MealRepository, foodRepo ports.FoodRepository) ports.MealService {
	return &mealService{
		mealRepo: mealRepo,
		foodRepo: foodRepo,
	}
}

func (s *mealService) GetMeals(ctx context.Context, userID string, date *time.Time) ([]*domain.Meal, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	meals, err := s.mealRepo.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	return meals, nil
}

func (s *mealService) GetMeal(ctx context.Context, mealID string) (*domain.Meal, error) {
	if mealID == "" {
		return nil, domain.ErrInvalidInput
	}

	meal, err := s.mealRepo.GetByID(ctx, mealID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}

	return meal, nil
}

func (s *mealService) CreateMeal(ctx context.Context, userID string, mealData *domain.Meal) (*domain.Meal, error) {
	if userID == "" || mealData == nil {
		return nil, domain.ErrInvalidInput
	}

	// Set defaults
	if mealData.ID == "" {
		mealData.ID = uuid.New().String()
	}
	mealData.UserID = userID

	if mealData.ConsumedAt.IsZero() {
		mealData.ConsumedAt = time.Now()
	}

	// Validate meal type
	validTypes := map[string]bool{
		"breakfast": true,
		"lunch":     true,
		"dinner":    true,
		"snack":     true,
	}
	if !validTypes[mealData.MealType] {
		return nil, domain.ErrInvalidInput
	}

	// Create meal
	if err := s.mealRepo.Create(ctx, mealData); err != nil {
		return nil, fmt.Errorf("failed to create meal: %w", err)
	}

	return mealData, nil
}

func (s *mealService) ConfirmParsedMeal(ctx context.Context, userID string, parsedMeal *domain.Meal) (*domain.Meal, error) {
	if userID == "" || parsedMeal == nil {
		return nil, domain.ErrInvalidInput
	}

	// Set defaults
	parsedMeal.ID = uuid.New().String()
	parsedMeal.UserID = userID
	parsedMeal.Source = "ai_parsed"

	if parsedMeal.ConsumedAt.IsZero() {
		parsedMeal.ConsumedAt = time.Now()
	}

	// Validate and process food items
	for i, item := range parsedMeal.Items {
		if item.FoodID == "" {
			return nil, domain.ErrInvalidInput
		}

		// Verify food exists
		food, err := s.foodRepo.GetByID(ctx, item.FoodID)
		if err != nil {
			return nil, fmt.Errorf("food item %d not found: %w", i, err)
		}

		// Set item ID if not present
		if item.ID == "" {
			item.ID = uuid.New().String()
		}
		item.MealID = parsedMeal.ID

		// Calculate nutrition based on quantity
		item.Calories = food.Calories * item.Quantity
		item.Protein = food.Protein * item.Quantity
		item.Carbs = food.Carbs * item.Quantity
		item.Fat = food.Fat * item.Quantity
	}

	// Create meal
	if err := s.mealRepo.Create(ctx, parsedMeal); err != nil {
		return nil, fmt.Errorf("failed to create parsed meal: %w", err)
	}

	return parsedMeal, nil
}

func (s *mealService) UpdateMeal(ctx context.Context, mealID string, updates map[string]interface{}) (*domain.Meal, error) {
	if mealID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify meal exists
	existing, err := s.mealRepo.GetByID(ctx, mealID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}

	// Validate meal type if being updated
	if mealType, ok := updates["meal_type"].(string); ok {
		validTypes := map[string]bool{
			"breakfast": true,
			"lunch":     true,
			"dinner":    true,
			"snack":     true,
		}
		if !validTypes[mealType] {
			return nil, domain.ErrInvalidInput
		}
	}

	// Update meal
	if err := s.mealRepo.Update(ctx, mealID, updates); err != nil {
		return nil, fmt.Errorf("failed to update meal: %w", err)
	}

	// Return updated meal
	return s.mealRepo.GetByID(ctx, existing.ID)
}

func (s *mealService) DeleteMeal(ctx context.Context, mealID string) error {
	if mealID == "" {
		return domain.ErrInvalidInput
	}

	// Verify meal exists
	_, err := s.mealRepo.GetByID(ctx, mealID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get meal: %w", err)
	}

	// Delete meal
	if err := s.mealRepo.Delete(ctx, mealID); err != nil {
		return fmt.Errorf("failed to delete meal: %w", err)
	}

	return nil
}

func (s *mealService) CalculateMealNutrition(ctx context.Context, mealID string) (*domain.NutritionTotals, error) {
	if mealID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Get meal with items
	meal, err := s.mealRepo.GetByID(ctx, mealID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}

	// Calculate totals
	totals := &domain.NutritionTotals{}
	for _, item := range meal.Items {
		totals.Calories += item.Calories
		totals.Protein += item.Protein
		totals.Carbs += item.Carbs
		totals.Fat += item.Fat
		totals.Fiber += item.Fiber
		totals.Sugar += item.Sugar
		totals.Sodium += item.Sodium
	}

	return totals, nil
}
