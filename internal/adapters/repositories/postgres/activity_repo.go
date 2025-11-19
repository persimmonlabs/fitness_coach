package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type activityRepository struct {
	db *gorm.DB
}

// NewActivityRepository creates a new activity repository
func NewActivityRepository(db *gorm.DB) ports.ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(ctx context.Context, activity *domain.Activity) error {
	return r.db.WithContext(ctx).Create(activity).Error
}

func (r *activityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Activity, error) {
	var activity domain.Activity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&activity).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) Update(ctx context.Context, activity *domain.Activity) error {
	return r.db.WithContext(ctx).Save(activity).Error
}

func (r *activityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Activity{}, "id = ?", id).Error
}

func (r *activityRepository) ListByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*domain.Activity, error) {
	var activities []*domain.Activity
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("start_time BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("start_time DESC").
		Find(&activities).Error

	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *activityRepository) GetTotalsByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	type Result struct {
		TotalCaloriesBurned float64
		TotalDuration       int
		TotalDistance       float64
		TotalSteps          int
	}

	var result Result
	query := r.db.WithContext(ctx).
		Model(&domain.Activity{}).
		Select(`
			COALESCE(SUM(calories_burned), 0) as total_calories_burned,
			COALESCE(SUM(duration_minutes), 0) as total_duration,
			COALESCE(SUM(distance), 0) as total_distance,
			COALESCE(SUM(steps), 0) as total_steps
		`).
		Where("user_id = ?", userID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("start_time BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_calories_burned": result.TotalCaloriesBurned,
		"total_duration":        result.TotalDuration,
		"total_distance":        result.TotalDistance,
		"total_steps":           result.TotalSteps,
	}, nil
}
