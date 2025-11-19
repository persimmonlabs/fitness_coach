package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type mealRepository struct {
	db *gorm.DB
}

// NewMealRepository creates a new meal repository
func NewMealRepository(db *gorm.DB) ports.MealRepository {
	return &mealRepository{db: db}
}

func (r *mealRepository) Create(ctx context.Context, meal *domain.Meal) error {
	return r.db.WithContext(ctx).Create(meal).Error
}

func (r *mealRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Meal, error) {
	var meal domain.Meal
	err := r.db.WithContext(ctx).
		Preload("FoodItems.Food").
		Where("id = ?", id).
		First(&meal).Error
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (r *mealRepository) Update(ctx context.Context, meal *domain.Meal) error {
	return r.db.WithContext(ctx).Save(meal).Error
}

func (r *mealRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Meal{}, "id = ?", id).Error
}

func (r *mealRepository) ListByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*domain.Meal, error) {
	var meals []*domain.Meal
	query := r.db.WithContext(ctx).
		Preload("FoodItems.Food").
		Where("user_id = ?", userID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("consumed_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("consumed_at DESC").
		Find(&meals).Error

	if err != nil {
		return nil, err
	}
	return meals, nil
}

// Food item operations

func (r *mealRepository) AddFoodItem(ctx context.Context, item *domain.MealFoodItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *mealRepository) UpdateFoodItem(ctx context.Context, item *domain.MealFoodItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *mealRepository) RemoveFoodItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.MealFoodItem{}, "id = ?", id).Error
}

func (r *mealRepository) GetFoodItems(ctx context.Context, mealID uuid.UUID) ([]*domain.MealFoodItem, error) {
	var items []*domain.MealFoodItem
	err := r.db.WithContext(ctx).
		Preload("Food").
		Where("meal_id = ?", mealID).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
