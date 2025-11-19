package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type goalRepository struct {
	db *gorm.DB
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *gorm.DB) ports.GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) Create(ctx context.Context, goal *domain.Goal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

func (r *goalRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Goal, error) {
	var goal domain.Goal
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&goal).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *goalRepository) Update(ctx context.Context, goal *domain.Goal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

func (r *goalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Goal{}, "id = ?", id).Error
}

func (r *goalRepository) ListByUser(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]*domain.Goal, error) {
	var goals []*domain.Goal
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&goals).Error

	if err != nil {
		return nil, err
	}
	return goals, nil
}
