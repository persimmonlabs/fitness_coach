package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type activityService struct {
	activityRepo ports.ActivityRepository
}

// NewActivityService creates a new activity service
func NewActivityService(activityRepo ports.ActivityRepository) ports.ActivityService {
	return &activityService{
		activityRepo: activityRepo,
	}
}

func (s *activityService) GetActivities(ctx context.Context, userID string, startDate, endDate *time.Time) ([]*domain.Activity, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	activities, err := s.activityRepo.GetByUserAndDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}

	return activities, nil
}

func (s *activityService) GetActivity(ctx context.Context, activityID string) (*domain.Activity, error) {
	if activityID == "" {
		return nil, domain.ErrInvalidInput
	}

	activity, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	return activity, nil
}

func (s *activityService) CreateActivity(ctx context.Context, userID string, activityData *domain.Activity) (*domain.Activity, error) {
	if userID == "" || activityData == nil {
		return nil, domain.ErrInvalidInput
	}

	// Validate required fields
	if activityData.ActivityType == "" {
		return nil, domain.ErrInvalidInput
	}

	// Set defaults
	if activityData.ID == "" {
		activityData.ID = uuid.New().String()
	}
	activityData.UserID = userID

	if activityData.StartTime.IsZero() {
		activityData.StartTime = time.Now()
	}

	// Validate activity type
	validTypes := map[string]bool{
		"walking":      true,
		"running":      true,
		"cycling":      true,
		"swimming":     true,
		"hiking":       true,
		"yoga":         true,
		"sports":       true,
		"other":        true,
	}
	if !validTypes[activityData.ActivityType] {
		return nil, domain.ErrInvalidInput
	}

	// Validate duration if provided
	if activityData.DurationMinutes < 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate calories burned if provided
	if activityData.CaloriesBurned < 0 {
		return nil, domain.ErrInvalidInput
	}

	// Create activity
	if err := s.activityRepo.Create(ctx, activityData); err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return activityData, nil
}

func (s *activityService) UpdateActivity(ctx context.Context, activityID string, updates map[string]interface{}) (*domain.Activity, error) {
	if activityID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify activity exists
	existing, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	// Validate activity type if being updated
	if activityType, ok := updates["activity_type"].(string); ok {
		validTypes := map[string]bool{
			"walking":      true,
			"running":      true,
			"cycling":      true,
			"swimming":     true,
			"hiking":       true,
			"yoga":         true,
			"sports":       true,
			"other":        true,
		}
		if !validTypes[activityType] {
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate duration if being updated
	if duration, ok := updates["duration_minutes"].(int); ok {
		if duration < 0 {
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate calories burned if being updated
	if calories, ok := updates["calories_burned"].(float64); ok {
		if calories < 0 {
			return nil, domain.ErrInvalidInput
		}
	}

	// Update activity
	if err := s.activityRepo.Update(ctx, activityID, updates); err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	// Return updated activity
	return s.activityRepo.GetByID(ctx, existing.ID)
}

func (s *activityService) DeleteActivity(ctx context.Context, activityID string) error {
	if activityID == "" {
		return domain.ErrInvalidInput
	}

	// Verify activity exists
	_, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get activity: %w", err)
	}

	// Delete activity
	if err := s.activityRepo.Delete(ctx, activityID); err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	return nil
}
