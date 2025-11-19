package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type goalService struct {
	goalRepo ports.GoalRepository
}

// NewGoalService creates a new goal service
func NewGoalService(goalRepo ports.GoalRepository) ports.GoalService {
	return &goalService{
		goalRepo: goalRepo,
	}
}

func (s *goalService) CreateGoal(ctx context.Context, userID string, goalData *domain.Goal) (*domain.Goal, error) {
	if userID == "" || goalData == nil {
		return nil, domain.ErrInvalidInput
	}

	// Validate required fields
	if goalData.GoalType == "" || goalData.TargetValue == 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate goal type
	validTypes := map[string]bool{
		"weight_loss":   true,
		"weight_gain":   true,
		"muscle_gain":   true,
		"fat_loss":      true,
		"endurance":     true,
		"strength":      true,
		"steps":         true,
		"calories":      true,
		"water_intake":  true,
		"sleep":         true,
		"other":         true,
	}
	if !validTypes[goalData.GoalType] {
		return nil, domain.ErrInvalidInput
	}

	// Set defaults
	if goalData.ID == "" {
		goalData.ID = uuid.New().String()
	}
	goalData.UserID = userID

	if goalData.Status == "" {
		goalData.Status = "active"
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":    true,
		"completed": true,
		"abandoned": true,
	}
	if !validStatuses[goalData.Status] {
		return nil, domain.ErrInvalidInput
	}

	// Validate target date if provided
	if !goalData.TargetDate.IsZero() && goalData.TargetDate.Before(time.Now()) {
		return nil, domain.ErrInvalidInput
	}

	// Create goal
	if err := s.goalRepo.Create(ctx, goalData); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	return goalData, nil
}

func (s *goalService) GetGoals(ctx context.Context, userID string, status *string) ([]*domain.Goal, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Validate status if provided
	if status != nil {
		validStatuses := map[string]bool{
			"active":    true,
			"completed": true,
			"abandoned": true,
		}
		if !validStatuses[*status] {
			return nil, domain.ErrInvalidInput
		}
	}

	goals, err := s.goalRepo.GetByUser(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals: %w", err)
	}

	return goals, nil
}

func (s *goalService) UpdateGoal(ctx context.Context, goalID string, updates map[string]interface{}) (*domain.Goal, error) {
	if goalID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify goal exists
	existing, err := s.goalRepo.GetByID(ctx, goalID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	// Validate status if being updated
	if status, ok := updates["status"].(string); ok {
		validStatuses := map[string]bool{
			"active":    true,
			"completed": true,
			"abandoned": true,
		}
		if !validStatuses[status] {
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate goal type if being updated
	if goalType, ok := updates["goal_type"].(string); ok {
		validTypes := map[string]bool{
			"weight_loss":   true,
			"weight_gain":   true,
			"muscle_gain":   true,
			"fat_loss":      true,
			"endurance":     true,
			"strength":      true,
			"steps":         true,
			"calories":      true,
			"water_intake":  true,
			"sleep":         true,
			"other":         true,
		}
		if !validTypes[goalType] {
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate target date if being updated
	if targetDate, ok := updates["target_date"].(time.Time); ok {
		if targetDate.Before(time.Now()) {
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate target value if being updated
	if targetValue, ok := updates["target_value"].(float64); ok {
		if targetValue <= 0 {
			return nil, domain.ErrInvalidInput
		}
	}

	// Update goal
	if err := s.goalRepo.Update(ctx, goalID, updates); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	// Return updated goal
	return s.goalRepo.GetByID(ctx, existing.ID)
}

func (s *goalService) DeleteGoal(ctx context.Context, goalID string) error {
	if goalID == "" {
		return domain.ErrInvalidInput
	}

	// Verify goal exists
	_, err := s.goalRepo.GetByID(ctx, goalID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get goal: %w", err)
	}

	// Delete goal
	if err := s.goalRepo.Delete(ctx, goalID); err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	return nil
}
