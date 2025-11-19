package integration

import (
	"testing"
	"time"

	"fitness-tracker/internal/adapters/repositories/postgres"
	"fitness-tracker/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateActivity(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "activity_test@example.com")

	// Initialize repository
	activityRepo := postgres.NewActivityRepository(testDB.DB)

	t.Run("Create cardio activity", func(t *testing.T) {
		duration := 45
		calories := 400.0
		distance := 5000.0 // 5km in meters
		avgHeartRate := 150
		maxHeartRate := 175

		activity := &domain.Activity{
			UserID:           user.ID,
			ActivityType:     "running",
			StartTime:        time.Now().Add(-1 * time.Hour),
			DurationMinutes:  &duration,
			CaloriesBurned:   &calories,
			Distance:         &distance,
			AverageHeartRate: &avgHeartRate,
			MaxHeartRate:     &maxHeartRate,
		}

		err := activityRepo.Create(activity)
		require.NoError(t, err)
		assert.NotEqual(t, "", activity.ID.String())

		// Retrieve and verify
		retrieved, err := activityRepo.GetByID(activity.ID)
		require.NoError(t, err)
		assert.Equal(t, "running", retrieved.ActivityType)
		assert.Equal(t, duration, *retrieved.DurationMinutes)
		assert.Equal(t, calories, *retrieved.CaloriesBurned)
		assert.Equal(t, distance, *retrieved.Distance)
		assert.Equal(t, avgHeartRate, *retrieved.AverageHeartRate)
	})

	t.Run("Create activity with notes", func(t *testing.T) {
		duration := 30
		calories := 250.0
		notes := "Great workout, felt energized"

		activity := &domain.Activity{
			UserID:          user.ID,
			ActivityType:    "cycling",
			StartTime:       time.Now(),
			DurationMinutes: &duration,
			CaloriesBurned:  &calories,
			Notes:           &notes,
		}

		err := activityRepo.Create(activity)
		require.NoError(t, err)

		// Retrieve and verify notes
		retrieved, err := activityRepo.GetByID(activity.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved.Notes)
		assert.Equal(t, notes, *retrieved.Notes)
	})
}

func TestGetActivitiesInRange(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "range_test@example.com")

	// Initialize repository
	activityRepo := postgres.NewActivityRepository(testDB.DB)

	t.Run("Get activities within date range", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		twoDaysAgo := now.Add(-48 * time.Hour)

		duration30 := 30
		duration45 := 45
		calories200 := 200.0
		calories300 := 300.0
		calories400 := 400.0

		// Create activities on different days
		activities := []*domain.Activity{
			{
				UserID:          user.ID,
				ActivityType:    "running",
				StartTime:       now,
				DurationMinutes: &duration45,
				CaloriesBurned:  &calories400,
			},
			{
				UserID:          user.ID,
				ActivityType:    "cycling",
				StartTime:       yesterday,
				DurationMinutes: &duration30,
				CaloriesBurned:  &calories300,
			},
			{
				UserID:          user.ID,
				ActivityType:    "swimming",
				StartTime:       twoDaysAgo,
				DurationMinutes: &duration30,
				CaloriesBurned:  &calories200,
			},
		}

		for _, activity := range activities {
			err := activityRepo.Create(activity)
			require.NoError(t, err)
		}

		// Get activities from yesterday onwards
		startDate := yesterday.Add(-1 * time.Hour)
		endDate := now.Add(1 * time.Hour)
		retrieved, err := activityRepo.GetByUserIDAndDateRange(user.ID, startDate, endDate)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2) // Should get today's and yesterday's activities
	})
}

func TestActivityUpdate(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "update_test@example.com")

	// Initialize repository
	activityRepo := postgres.NewActivityRepository(testDB.DB)

	t.Run("Update activity details", func(t *testing.T) {
		// Create activity
		activity := CreateTestActivity(t, testDB.DB, user.ID, "walking")

		// Update activity
		newDuration := 60
		newCalories := 350.0
		newNotes := "Updated workout notes"

		activity.DurationMinutes = &newDuration
		activity.CaloriesBurned = &newCalories
		activity.Notes = &newNotes

		err := activityRepo.Update(activity)
		require.NoError(t, err)

		// Retrieve and verify
		updated, err := activityRepo.GetByID(activity.ID)
		require.NoError(t, err)
		assert.Equal(t, newDuration, *updated.DurationMinutes)
		assert.Equal(t, newCalories, *updated.CaloriesBurned)
		assert.NotNil(t, updated.Notes)
		assert.Equal(t, newNotes, *updated.Notes)
	})
}

func TestActivityDelete(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "delete_test@example.com")

	// Initialize repository
	activityRepo := postgres.NewActivityRepository(testDB.DB)

	t.Run("Delete activity", func(t *testing.T) {
		// Create activity
		activity := CreateTestActivity(t, testDB.DB, user.ID, "yoga")

		// Delete activity
		err := activityRepo.Delete(activity.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = activityRepo.GetByID(activity.ID)
		assert.Error(t, err)
	})
}

func TestGetActivitiesByType(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "type_test@example.com")

	// Initialize repository
	activityRepo := postgres.NewActivityRepository(testDB.DB)

	t.Run("Filter activities by type", func(t *testing.T) {
		// Create activities of different types
		CreateTestActivity(t, testDB.DB, user.ID, "running")
		CreateTestActivity(t, testDB.DB, user.ID, "running")
		CreateTestActivity(t, testDB.DB, user.ID, "cycling")
		CreateTestActivity(t, testDB.DB, user.ID, "swimming")

		// Query running activities
		var runningActivities []domain.Activity
		err := testDB.DB.Where("user_id = ? AND activity_type = ?", user.ID, "running").
			Find(&runningActivities).Error
		require.NoError(t, err)
		assert.Len(t, runningActivities, 2)

		// Verify all are running activities
		for _, activity := range runningActivities {
			assert.Equal(t, "running", activity.ActivityType)
		}
	})
}

func TestActivityStatistics(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "stats_test@example.com")

	t.Run("Calculate total calories burned", func(t *testing.T) {
		// Create multiple activities
		duration30 := 30
		calories200 := 200.0
		calories300 := 300.0
		calories400 := 400.0

		activities := []*domain.Activity{
			{
				UserID:          user.ID,
				ActivityType:    "running",
				StartTime:       time.Now(),
				DurationMinutes: &duration30,
				CaloriesBurned:  &calories400,
			},
			{
				UserID:          user.ID,
				ActivityType:    "cycling",
				StartTime:       time.Now(),
				DurationMinutes: &duration30,
				CaloriesBurned:  &calories300,
			},
			{
				UserID:          user.ID,
				ActivityType:    "swimming",
				StartTime:       time.Now(),
				DurationMinutes: &duration30,
				CaloriesBurned:  &calories200,
			},
		}

		for _, activity := range activities {
			err := testDB.DB.Create(activity).Error
			require.NoError(t, err)
		}

		// Calculate total calories
		var totalCalories float64
		err := testDB.DB.Model(&domain.Activity{}).
			Where("user_id = ?", user.ID).
			Select("COALESCE(SUM(calories_burned), 0)").
			Scan(&totalCalories).Error
		require.NoError(t, err)
		assert.Equal(t, 900.0, totalCalories)
	})
}
