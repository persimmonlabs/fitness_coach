package integration

import (
	"math"
	"testing"
	"time"

	"fitness-tracker/internal/adapters/repositories/postgres"
	"fitness-tracker/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkoutFlow(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "workout_test@example.com")

	// Create test exercise
	exercise := CreateTestExercise(t, testDB.DB, "Bench Press", "strength")

	// Initialize repository
	workoutRepo := postgres.NewWorkoutRepository(testDB.DB)

	t.Run("Complete workout flow", func(t *testing.T) {
		// Start workout
		startTime := time.Now()
		workout := &domain.Workout{
			UserID:    user.ID,
			Name:      "Chest Day",
			StartTime: startTime,
		}

		err := workoutRepo.Create(workout)
		require.NoError(t, err)
		assert.NotEqual(t, "", workout.ID.String())

		// Add exercise to workout
		workoutExercise := &domain.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exercise.ID,
			OrderIndex: 1,
		}

		err = testDB.DB.Create(workoutExercise).Error
		require.NoError(t, err)

		// Log sets
		sets := []domain.WorkoutSet{
			{
				WorkoutExerciseID: workoutExercise.ID,
				SetNumber:         1,
				Reps:              intPtr(10),
				Weight:            float64Ptr(60.0),
				RestSeconds:       intPtr(90),
			},
			{
				WorkoutExerciseID: workoutExercise.ID,
				SetNumber:         2,
				Reps:              intPtr(8),
				Weight:            float64Ptr(65.0),
				RestSeconds:       intPtr(90),
			},
			{
				WorkoutExerciseID: workoutExercise.ID,
				SetNumber:         3,
				Reps:              intPtr(6),
				Weight:            float64Ptr(70.0),
				RestSeconds:       intPtr(120),
			},
		}

		for _, set := range sets {
			err = testDB.DB.Create(&set).Error
			require.NoError(t, err)
		}

		// Finish workout
		endTime := time.Now()
		durationMinutes := int(endTime.Sub(startTime).Minutes())
		caloriesBurned := 250.0

		workout.EndTime = &endTime
		workout.DurationMinutes = &durationMinutes
		workout.CaloriesBurned = &caloriesBurned

		err = workoutRepo.Update(workout)
		require.NoError(t, err)

		// Retrieve complete workout
		retrieved, err := workoutRepo.GetByID(workout.ID)
		require.NoError(t, err)
		assert.Equal(t, "Chest Day", retrieved.Name)
		assert.NotNil(t, retrieved.EndTime)
		assert.NotNil(t, retrieved.DurationMinutes)

		// Load exercises and sets
		err = testDB.DB.Preload("Exercise").Find(&retrieved.Exercises, "workout_id = ?", workout.ID).Error
		require.NoError(t, err)
		assert.Len(t, retrieved.Exercises, 1)

		// Load sets for the exercise
		err = testDB.DB.Find(&retrieved.Exercises[0].Sets, "workout_exercise_id = ?", workoutExercise.ID).Error
		require.NoError(t, err)
		assert.Len(t, retrieved.Exercises[0].Sets, 3)

		// Verify estimated 1RM calculation (using Brzycki formula)
		// 1RM = weight × (36 / (37 - reps))
		heaviestSet := sets[2] // 70kg x 6 reps
		estimated1RM := *heaviestSet.Weight * (36.0 / (37.0 - float64(*heaviestSet.Reps)))
		assert.InDelta(t, 81.29, estimated1RM, 0.1, "Estimated 1RM should be around 81.29kg")
	})
}

func TestWorkoutWithMultipleExercises(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "multi_exercise_test@example.com")

	// Create test exercises
	exercise1 := CreateTestExercise(t, testDB.DB, "Squat", "strength")
	exercise2 := CreateTestExercise(t, testDB.DB, "Deadlift", "strength")
	exercise3 := CreateTestExercise(t, testDB.DB, "Leg Press", "strength")

	// Initialize repository
	workoutRepo := postgres.NewWorkoutRepository(testDB.DB)

	t.Run("Workout with multiple exercises", func(t *testing.T) {
		// Create workout
		workout := &domain.Workout{
			UserID:    user.ID,
			Name:      "Leg Day",
			StartTime: time.Now(),
		}

		err := workoutRepo.Create(workout)
		require.NoError(t, err)

		// Add exercises in order
		exercises := []*domain.WorkoutExercise{
			{WorkoutID: workout.ID, ExerciseID: exercise1.ID, OrderIndex: 1},
			{WorkoutID: workout.ID, ExerciseID: exercise2.ID, OrderIndex: 2},
			{WorkoutID: workout.ID, ExerciseID: exercise3.ID, OrderIndex: 3},
		}

		for _, ex := range exercises {
			err = testDB.DB.Create(ex).Error
			require.NoError(t, err)
		}

		// Retrieve and verify order
		var retrievedExercises []domain.WorkoutExercise
		err = testDB.DB.Where("workout_id = ?", workout.ID).
			Order("order_index ASC").
			Find(&retrievedExercises).Error
		require.NoError(t, err)
		assert.Len(t, retrievedExercises, 3)
		assert.Equal(t, 1, retrievedExercises[0].OrderIndex)
		assert.Equal(t, 2, retrievedExercises[1].OrderIndex)
		assert.Equal(t, 3, retrievedExercises[2].OrderIndex)
	})
}

func TestWorkoutDelete(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "delete_workout_test@example.com")

	// Initialize repository
	workoutRepo := postgres.NewWorkoutRepository(testDB.DB)

	t.Run("Delete workout cascades to exercises and sets", func(t *testing.T) {
		// Create workout with exercise and sets
		exercise := CreateTestExercise(t, testDB.DB, "Push-ups", "bodyweight")

		workout := &domain.Workout{
			UserID:    user.ID,
			Name:      "Quick Workout",
			StartTime: time.Now(),
		}
		err := workoutRepo.Create(workout)
		require.NoError(t, err)

		workoutExercise := &domain.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exercise.ID,
			OrderIndex: 1,
		}
		err = testDB.DB.Create(workoutExercise).Error
		require.NoError(t, err)

		set := &domain.WorkoutSet{
			WorkoutExerciseID: workoutExercise.ID,
			SetNumber:         1,
			Reps:              intPtr(20),
		}
		err = testDB.DB.Create(set).Error
		require.NoError(t, err)

		// Delete workout
		err = workoutRepo.Delete(workout.ID)
		require.NoError(t, err)

		// Verify workout is deleted
		_, err = workoutRepo.GetByID(workout.ID)
		assert.Error(t, err)

		// Verify exercises and sets are also deleted (due to cascade)
		var exerciseCount int64
		testDB.DB.Model(&domain.WorkoutExercise{}).Where("workout_id = ?", workout.ID).Count(&exerciseCount)
		assert.Equal(t, int64(0), exerciseCount)

		var setCount int64
		testDB.DB.Model(&domain.WorkoutSet{}).Where("workout_exercise_id = ?", workoutExercise.ID).Count(&setCount)
		assert.Equal(t, int64(0), setCount)
	})
}

func TestGetWorkoutsInRange(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "range_workout_test@example.com")

	// Initialize repository
	workoutRepo := postgres.NewWorkoutRepository(testDB.DB)

	t.Run("Get workouts within date range", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		twoDaysAgo := now.Add(-48 * time.Hour)

		// Create workouts on different days
		workouts := []*domain.Workout{
			{UserID: user.ID, Name: "Workout Today", StartTime: now},
			{UserID: user.ID, Name: "Workout Yesterday", StartTime: yesterday},
			{UserID: user.ID, Name: "Workout Two Days Ago", StartTime: twoDaysAgo},
		}

		for _, workout := range workouts {
			err := workoutRepo.Create(workout)
			require.NoError(t, err)
		}

		// Get workouts from yesterday onwards
		startDate := yesterday.Add(-1 * time.Hour)
		endDate := now.Add(1 * time.Hour)
		retrieved, err := workoutRepo.GetByUserIDAndDateRange(user.ID, startDate, endDate)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2) // Should get today's and yesterday's workouts
	})
}

func TestWorkoutWithCardioExercise(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "cardio_test@example.com")

	// Create cardio exercise
	exercise := CreateTestExercise(t, testDB.DB, "Running", "cardio")

	// Initialize repository
	workoutRepo := postgres.NewWorkoutRepository(testDB.DB)

	t.Run("Log cardio workout with duration and distance", func(t *testing.T) {
		workout := &domain.Workout{
			UserID:    user.ID,
			Name:      "Morning Run",
			StartTime: time.Now(),
		}

		err := workoutRepo.Create(workout)
		require.NoError(t, err)

		workoutExercise := &domain.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exercise.ID,
			OrderIndex: 1,
		}
		err = testDB.DB.Create(workoutExercise).Error
		require.NoError(t, err)

		// Log cardio set with duration and distance
		set := &domain.WorkoutSet{
			WorkoutExerciseID: workoutExercise.ID,
			SetNumber:         1,
			DurationSeconds:   intPtr(1800), // 30 minutes
			Distance:          float64Ptr(5000.0), // 5km in meters
		}

		err = testDB.DB.Create(set).Error
		require.NoError(t, err)

		// Retrieve and verify
		var retrievedSet domain.WorkoutSet
		err = testDB.DB.First(&retrievedSet, "id = ?", set.ID).Error
		require.NoError(t, err)
		assert.Equal(t, 1800, *retrievedSet.DurationSeconds)
		assert.Equal(t, 5000.0, *retrievedSet.Distance)

		// Calculate pace (min/km)
		durationMinutes := float64(*retrievedSet.DurationSeconds) / 60.0
		distanceKm := *retrievedSet.Distance / 1000.0
		pace := durationMinutes / distanceKm
		assert.InDelta(t, 6.0, pace, 0.1, "Pace should be 6 min/km")
	})
}

func TestEstimated1RMCalculation(t *testing.T) {
	t.Run("Calculate 1RM using Brzycki formula", func(t *testing.T) {
		testCases := []struct {
			weight      float64
			reps        int
			expected1RM float64
		}{
			{60.0, 10, 80.0},   // 60kg x 10 reps ≈ 80kg 1RM
			{100.0, 5, 112.5},  // 100kg x 5 reps ≈ 112.5kg 1RM
			{80.0, 8, 100.0},   // 80kg x 8 reps ≈ 100kg 1RM
		}

		for _, tc := range testCases {
			// Brzycki formula: 1RM = weight × (36 / (37 - reps))
			estimated1RM := tc.weight * (36.0 / (37.0 - float64(tc.reps)))
			assert.InDelta(t, tc.expected1RM, estimated1RM, 2.5,
				"1RM calculation for %vkg x %v reps", tc.weight, tc.reps)
		}
	})
}
