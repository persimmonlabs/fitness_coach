package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type foodRepository struct {
	db *gorm.DB
}

// NewFoodRepository creates a new food repository
func NewFoodRepository(db *gorm.DB) ports.FoodRepository {
	return &foodRepository{db: db}
}

func (r *foodRepository) Create(ctx context.Context, food *domain.Food) error {
	return r.db.WithContext(ctx).Create(food).Error
}

func (r *foodRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Food, error) {
	var food domain.Food
	err := r.db.WithContext(ctx).
		Preload("Ingredients.Ingredient").
		Preload("ServingConversions.ServingUnit").
		Where("id = ?", id).
		First(&food).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) GetByFdcID(ctx context.Context, fdcID int) (*domain.Food, error) {
	var food domain.Food
	err := r.db.WithContext(ctx).
		Preload("Ingredients.Ingredient").
		Preload("ServingConversions.ServingUnit").
		Where("fdc_id = ?", fdcID).
		First(&food).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) Update(ctx context.Context, food *domain.Food) error {
	return r.db.WithContext(ctx).Save(food).Error
}

func (r *foodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Food{}, "id = ?", id).Error
}

func (r *foodRepository) List(ctx context.Context, limit, offset int) ([]*domain.Food, error) {
	var foods []*domain.Food
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&foods).Error
	if err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *foodRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Food, error) {
	var foods []*domain.Food

	// Use PostgreSQL full-text search with to_tsquery
	err := r.db.WithContext(ctx).
		Where("to_tsvector('english', name) @@ plainto_tsquery('english', ?)", query).
		Or("name ILIKE ?", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Order("is_verified DESC, name ASC").
		Find(&foods).Error

	if err != nil {
		return nil, err
	}
	return foods, nil
}

// Ingredient operations

func (r *foodRepository) AddIngredient(ctx context.Context, ingredient *domain.FoodIngredient) error {
	return r.db.WithContext(ctx).Create(ingredient).Error
}

func (r *foodRepository) GetIngredients(ctx context.Context, foodID uuid.UUID) ([]*domain.FoodIngredient, error) {
	var ingredients []*domain.FoodIngredient
	err := r.db.WithContext(ctx).
		Preload("Ingredient").
		Where("food_id = ?", foodID).
		Find(&ingredients).Error
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

// Serving conversion operations

func (r *foodRepository) AddServingConversion(ctx context.Context, conversion *domain.FoodServingConversion) error {
	return r.db.WithContext(ctx).Create(conversion).Error
}

func (r *foodRepository) GetServingConversions(ctx context.Context, foodID uuid.UUID) ([]*domain.FoodServingConversion, error) {
	var conversions []*domain.FoodServingConversion
	err := r.db.WithContext(ctx).
		Preload("ServingUnit").
		Where("food_id = ?", foodID).
		Find(&conversions).Error
	if err != nil {
		return nil, err
	}
	return conversions, nil
}

// Serving unit operations

func (r *foodRepository) CreateServingUnit(ctx context.Context, unit *domain.ServingUnit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *foodRepository) GetServingUnit(ctx context.Context, id uuid.UUID) (*domain.ServingUnit, error) {
	var unit domain.ServingUnit
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *foodRepository) ListServingUnits(ctx context.Context) ([]*domain.ServingUnit, error) {
	var units []*domain.ServingUnit
	err := r.db.WithContext(ctx).
		Order("category ASC, name ASC").
		Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}
