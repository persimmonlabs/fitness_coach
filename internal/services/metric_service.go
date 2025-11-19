package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type metricService struct {
	metricRepo ports.MetricRepository
}

// NewMetricService creates a new metric service
func NewMetricService(metricRepo ports.MetricRepository) ports.MetricService {
	return &metricService{
		metricRepo: metricRepo,
	}
}

func (s *metricService) LogMetric(ctx context.Context, userID, metricType string, value float64, unit string, recordedAt time.Time) (*domain.Metric, error) {
	if userID == "" || metricType == "" || unit == "" {
		return nil, domain.ErrInvalidInput
	}

	// Validate metric type
	validTypes := map[string]bool{
		"weight":         true,
		"body_fat":       true,
		"muscle_mass":    true,
		"water":          true,
		"bmi":            true,
		"blood_pressure": true,
		"heart_rate":     true,
		"steps":          true,
		"sleep":          true,
		"other":          true,
	}
	if !validTypes[metricType] {
		return nil, domain.ErrInvalidInput
	}

	// Validate value
	if value < 0 {
		return nil, domain.ErrInvalidInput
	}

	// Set defaults
	metric := &domain.Metric{
		ID:         uuid.New().String(),
		UserID:     userID,
		MetricType: metricType,
		Value:      value,
		Unit:       unit,
		RecordedAt: recordedAt,
	}

	if metric.RecordedAt.IsZero() {
		metric.RecordedAt = time.Now()
	}

	// Create metric
	if err := s.metricRepo.Create(ctx, metric); err != nil {
		return nil, fmt.Errorf("failed to log metric: %w", err)
	}

	return metric, nil
}

func (s *metricService) GetMetricTrend(ctx context.Context, userID, metricType string, startDate, endDate *time.Time) ([]*domain.Metric, error) {
	if userID == "" || metricType == "" {
		return nil, domain.ErrInvalidInput
	}

	// Validate metric type
	validTypes := map[string]bool{
		"weight":         true,
		"body_fat":       true,
		"muscle_mass":    true,
		"water":          true,
		"bmi":            true,
		"blood_pressure": true,
		"heart_rate":     true,
		"steps":          true,
		"sleep":          true,
		"other":          true,
	}
	if !validTypes[metricType] {
		return nil, domain.ErrInvalidInput
	}

	metrics, err := s.metricRepo.GetByUserAndType(ctx, userID, metricType, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric trend: %w", err)
	}

	return metrics, nil
}

func (s *metricService) GetLatestMetric(ctx context.Context, userID, metricType string) (*domain.Metric, error) {
	if userID == "" || metricType == "" {
		return nil, domain.ErrInvalidInput
	}

	// Validate metric type
	validTypes := map[string]bool{
		"weight":         true,
		"body_fat":       true,
		"muscle_mass":    true,
		"water":          true,
		"bmi":            true,
		"blood_pressure": true,
		"heart_rate":     true,
		"steps":          true,
		"sleep":          true,
		"other":          true,
	}
	if !validTypes[metricType] {
		return nil, domain.ErrInvalidInput
	}

	metric, err := s.metricRepo.GetLatest(ctx, userID, metricType)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest metric: %w", err)
	}

	return metric, nil
}
