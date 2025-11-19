package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type workoutService struct {
	workoutRepo  ports.WorkoutRepository
	exerciseRepo ports.ExerciseRepository
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(workoutRepo ports.WorkoutRepository, exerciseRepo ports.ExerciseRepository) ports.WorkoutService {
	return &workoutService{
		workoutRepo:  workoutRepo,
		exerciseRepo: exerciseRepo,
	}
}

func (s *workoutService) StartWorkout(ctx context.Context, userID, name string) (*domain.Workout, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	workout := &domain.Workout{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		Status:    "in_progress",
		StartTime: time.Now(),
	}

	if err := s.workoutRepo.Create(ctx, workout); err != nil {
		return nil, fmt.Errorf("failed to start workout: %w", err)
	}

	return workout, nil
}

func (s *workoutService) GetWorkouts(ctx context.Context, userID string, startDate, endDate *time.Time) ([]*domain.Workout, error) {
	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	workouts, err := s.workoutRepo.GetByUserAndDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get workouts: %w", err)
	}

	return workouts, nil
}

func (s *workoutService) GetWorkout(ctx context.Context, workoutID string) (*domain.Workout, error) {
	if workoutID == "" {
		return nil, domain.ErrInvalidInput
	}

	workout, err := s.workoutRepo.GetByID(ctx, workoutID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}

	return workout, nil
}

func (s *workoutService) AddExercise(ctx context.Context, workoutID, exerciseID string) (*domain.WorkoutExercise, error) {
	if workoutID == "" || exerciseID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify workout exists and is in progress
	workout, err := s.workoutRepo.GetByID(ctx, workoutID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}

	if workout.Status != "in_progress" {
		return nil, domain.ErrInvalidInput
	}

	// Verify exercise exists
	exercise, err := s.exerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	// Create workout exercise
	workoutExercise := &domain.WorkoutExercise{
		ID:         uuid.New().String(),
		WorkoutID:  workoutID,
		ExerciseID: exerciseID,
		Sets:       []domain.WorkoutSet{},
	}

	// Note: Repository should handle adding to workout.Exercises
	if err := s.workoutRepo.AddExercise(ctx, workoutExercise); err != nil {
		return nil, fmt.Errorf("failed to add exercise to workout: %w", err)
	}

	workoutExercise.Exercise = exercise
	return workoutExercise, nil
}

func (s *workoutService) LogSet(ctx context.Context, workoutExerciseID string, setData *domain.WorkoutSet) (*domain.WorkoutSet, error) {
	if workoutExerciseID == "" || setData == nil {
		return nil, domain.ErrInvalidInput
	}

	// Validate set data
	if setData.Reps < 0 {
		return nil, domain.ErrInvalidInput
	}
	if setData.Weight != nil && *setData.Weight < 0 {
		return nil, domain.ErrInvalidInput
	}
	if setData.Duration != nil && *setData.Duration < 0 {
		return nil, domain.ErrInvalidInput
	}

	// Set defaults
	if setData.ID == "" {
		setData.ID = uuid.New().String()
	}
	setData.WorkoutExerciseID = workoutExerciseID

	// Create set
	if err := s.workoutRepo.LogSet(ctx, setData); err != nil {
		return nil, fmt.Errorf("failed to log set: %w", err)
	}

	return setData, nil
}

func (s *workoutService) FinishWorkout(ctx context.Context, workoutID string) error {
	if workoutID == "" {
		return domain.ErrInvalidInput
	}

	// Verify workout exists and is in progress
	workout, err := s.workoutRepo.GetByID(ctx, workoutID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get workout: %w", err)
	}

	if workout.Status != "in_progress" {
		return domain.ErrInvalidInput
	}

	// Update workout status and end time
	updates := map[string]interface{}{
		"status":   "completed",
		"end_time": time.Now(),
	}

	if err := s.workoutRepo.Update(ctx, workoutID, updates); err != nil {
		return fmt.Errorf("failed to finish workout: %w", err)
	}

	return nil
}

func (s *workoutService) DeleteWorkout(ctx context.Context, workoutID string) error {
	if workoutID == "" {
		return domain.ErrInvalidInput
	}

	// Verify workout exists
	_, err := s.workoutRepo.GetByID(ctx, workoutID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get workout: %w", err)
	}

	// Delete workout (cascade should handle exercises and sets)
	if err := s.workoutRepo.Delete(ctx, workoutID); err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	return nil
}
