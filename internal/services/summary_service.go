package services

import (
	"context"
	"fmt"
	"time"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type summaryService struct {
	mealRepo     ports.MealRepository
	activityRepo ports.ActivityRepository
	workoutRepo  ports.WorkoutRepository
}

// NewSummaryService creates a new summary service
func NewSummaryService(
	mealRepo ports.MealRepository,
	activityRepo ports.ActivityRepository,
	workoutRepo ports.WorkoutRepository,
) ports.SummaryService {
	return &summaryService{
		mealRepo:     mealRepo,
		activityRepo: activityRepo,
		workoutRepo:  workoutRepo,
	}
}

func (s *summaryService) GetDailySummary(ctx context.Context, userID string, date time.Time) (*domain.DailySummary, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Calculate summary
	return s.CalculateDailySummary(ctx, userID, date)
}

func (s *summaryService) CalculateDailySummary(ctx context.Context, userID string, date time.Time) (*domain.DailySummary, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Normalize date to start of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	summary := &domain.DailySummary{
		UserID: userID,
		Date:   startOfDay,
		Nutrition: domain.NutritionTotals{
			Calories: 0,
			Protein:  0,
			Carbs:    0,
			Fat:      0,
			Fiber:    0,
			Sugar:    0,
			Sodium:   0,
		},
		CaloriesBurned:  0,
		WorkoutsCount:   0,
		ActivitiesCount: 0,
		MealsCount:      0,
	}

	// Get meals for the day
	meals, err := s.mealRepo.GetByUserAndDate(ctx, userID, &startOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	// Calculate nutrition totals from meals
	summary.MealsCount = len(meals)
	for _, meal := range meals {
		for _, item := range meal.Items {
			summary.Nutrition.Calories += item.Calories
			summary.Nutrition.Protein += item.Protein
			summary.Nutrition.Carbs += item.Carbs
			summary.Nutrition.Fat += item.Fat
			summary.Nutrition.Fiber += item.Fiber
			summary.Nutrition.Sugar += item.Sugar
			summary.Nutrition.Sodium += item.Sodium
		}
	}

	// Get activities for the day
	activities, err := s.activityRepo.GetByUserAndDateRange(ctx, userID, &startOfDay, &endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}

	// Calculate calories burned from activities
	summary.ActivitiesCount = len(activities)
	for _, activity := range activities {
		summary.CaloriesBurned += activity.CaloriesBurned
	}

	// Get workouts for the day
	workouts, err := s.workoutRepo.GetByUserAndDateRange(ctx, userID, &startOfDay, &endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get workouts: %w", err)
	}

	summary.WorkoutsCount = len(workouts)

	// Note: Calories burned from workouts could be calculated based on exercises
	// For now, we'll focus on activities for calorie tracking

	return summary, nil
}
