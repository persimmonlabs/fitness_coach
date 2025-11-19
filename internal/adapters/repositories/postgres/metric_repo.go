package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type metricRepository struct {
	db *gorm.DB
}

// NewMetricRepository creates a new metric repository
func NewMetricRepository(db *gorm.DB) ports.MetricRepository {
	return &metricRepository{db: db}
}

func (r *metricRepository) Create(ctx context.Context, metric *domain.Metric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

func (r *metricRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Metric, error) {
	var metric domain.Metric
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *metricRepository) Update(ctx context.Context, metric *domain.Metric) error {
	return r.db.WithContext(ctx).Save(metric).Error
}

func (r *metricRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Metric{}, "id = ?", id).Error
}

func (r *metricRepository) ListByUser(ctx context.Context, userID uuid.UUID, metricType string, startDate, endDate time.Time, limit, offset int) ([]*domain.Metric, error) {
	var metrics []*domain.Metric
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("measured_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("measured_at DESC").
		Find(&metrics).Error

	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// Daily summary operations

func (r *metricRepository) CreateOrUpdateDailySummary(ctx context.Context, summary *domain.DailySummary) error {
	// Use GORM's upsert functionality (INSERT ... ON CONFLICT)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"total_calories",
				"total_protein",
				"total_carbohydrates",
				"total_fat",
				"total_calories_burned",
				"total_exercise_minutes",
				"total_steps",
				"total_distance",
				"weight",
				"body_fat",
				"updated_at",
			}),
		}).
		Create(summary).Error
}

func (r *metricRepository) GetDailySummary(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.DailySummary, error) {
	var summary domain.DailySummary
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		First(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *metricRepository) ListDailySummaries(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*domain.DailySummary, error) {
	var summaries []*domain.DailySummary
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("date BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	err := query.
		Order("date DESC").
		Find(&summaries).Error

	if err != nil {
		return nil, err
	}
	return summaries, nil
}
