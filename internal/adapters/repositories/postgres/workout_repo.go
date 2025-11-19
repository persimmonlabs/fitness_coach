package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type workoutRepository struct {
	db *gorm.DB
}

// NewWorkoutRepository creates a new workout repository
func NewWorkoutRepository(db *gorm.DB) ports.WorkoutRepository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) Create(ctx context.Context, workout *domain.Workout) error {
	return r.db.WithContext(ctx).Create(workout).Error
}

func (r *workoutRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workout, error) {
	var workout domain.Workout
	err := r.db.WithContext(ctx).
		Preload("Exercises.Exercise").
		Preload("Exercises.Sets").
		Where("id = ?", id).
		First(&workout).Error
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (r *workoutRepository) Update(ctx context.Context, workout *domain.Workout) error {
	return r.db.WithContext(ctx).Save(workout).Error
}

func (r *workoutRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Workout{}, "id = ?", id).Error
}

func (r *workoutRepository) ListByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*domain.Workout, error) {
	var workouts []*domain.Workout
	query := r.db.WithContext(ctx).
		Preload("Exercises.Exercise").
		Preload("Exercises.Sets").
		Where("user_id = ?", userID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("start_time BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("start_time DESC").
		Find(&workouts).Error

	if err != nil {
		return nil, err
	}
	return workouts, nil
}

// Exercise operations

func (r *workoutRepository) CreateExercise(ctx context.Context, exercise *domain.Exercise) error {
	return r.db.WithContext(ctx).Create(exercise).Error
}

func (r *workoutRepository) GetExercise(ctx context.Context, id uuid.UUID) (*domain.Exercise, error) {
	var exercise domain.Exercise
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&exercise).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *workoutRepository) ListExercises(ctx context.Context, category string, limit, offset int) ([]*domain.Exercise, error) {
	var exercises []*domain.Exercise
	query := r.db.WithContext(ctx)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&exercises).Error

	if err != nil {
		return nil, err
	}
	return exercises, nil
}

// Workout exercise operations

func (r *workoutRepository) AddWorkoutExercise(ctx context.Context, workoutExercise *domain.WorkoutExercise) error {
	return r.db.WithContext(ctx).Create(workoutExercise).Error
}

func (r *workoutRepository) GetWorkoutExercises(ctx context.Context, workoutID uuid.UUID) ([]*domain.WorkoutExercise, error) {
	var workoutExercises []*domain.WorkoutExercise
	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Preload("Sets").
		Where("workout_id = ?", workoutID).
		Order("order_index ASC").
		Find(&workoutExercises).Error
	if err != nil {
		return nil, err
	}
	return workoutExercises, nil
}

// Set operations

func (r *workoutRepository) AddSet(ctx context.Context, set *domain.WorkoutSet) error {
	return r.db.WithContext(ctx).Create(set).Error
}

func (r *workoutRepository) UpdateSet(ctx context.Context, set *domain.WorkoutSet) error {
	return r.db.WithContext(ctx).Save(set).Error
}

func (r *workoutRepository) DeleteSet(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.WorkoutSet{}, "id = ?", id).Error
}

func (r *workoutRepository) GetSets(ctx context.Context, workoutExerciseID uuid.UUID) ([]*domain.WorkoutSet, error) {
	var sets []*domain.WorkoutSet
	err := r.db.WithContext(ctx).
		Where("workout_exercise_id = ?", workoutExerciseID).
		Order("set_number ASC").
		Find(&sets).Error
	if err != nil {
		return nil, err
	}
	return sets, nil
}
