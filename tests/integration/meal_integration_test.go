package integration

import (
	"testing"
	"time"

	"fitness-tracker/internal/adapters/repositories/postgres"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMealFlow(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "meal_test@example.com")

	// Create test food
	food := CreateTestFood(t, testDB.DB, "Chicken Breast", 165.0)

	// Initialize repositories and services
	mealRepo := postgres.NewMealRepository(testDB.DB)
	foodRepo := postgres.NewFoodRepository(testDB.DB)
	mealService := services.NewMealService(mealRepo, foodRepo)

	t.Run("Create and confirm meal", func(t *testing.T) {
		// Create meal
		meal := &domain.Meal{
			UserID:             user.ID,
			Name:               "Lunch",
			MealType:           "lunch",
			ConsumedAt:         time.Now(),
			TotalCalories:      0,
			TotalProtein:       0,
			TotalCarbohydrates: 0,
			TotalFat:           0,
		}

		err := mealRepo.Create(meal)
		require.NoError(t, err)
		assert.NotEqual(t, "", meal.ID.String())

		// Add food item
		foodItem := &domain.MealFoodItem{
			MealID:        meal.ID,
			FoodID:        food.ID,
			Quantity:      200, // 200g
			Unit:          "g",
			Calories:      330.0, // 165 * 2
			Protein:       20.0,
			Carbohydrates: 0.0,
			Fat:           6.6,
		}

		err = testDB.DB.Create(foodItem).Error
		require.NoError(t, err)

		// Retrieve meal with food items
		retrievedMeal, err := mealRepo.GetByID(meal.ID)
		require.NoError(t, err)
		assert.Equal(t, meal.ID, retrievedMeal.ID)
		assert.Equal(t, "Lunch", retrievedMeal.Name)

		// Load food items
		err = testDB.DB.Preload("Food").Find(&retrievedMeal.FoodItems, "meal_id = ?", meal.ID).Error
		require.NoError(t, err)
		assert.Len(t, retrievedMeal.FoodItems, 1)
		assert.Equal(t, food.ID, retrievedMeal.FoodItems[0].FoodID)

		// Update meal totals
		meal.TotalCalories = 330.0
		meal.TotalProtein = 20.0
		meal.TotalCarbohydrates = 0.0
		meal.TotalFat = 6.6
		err = mealRepo.Update(meal)
		require.NoError(t, err)

		// Verify nutrition calculation
		updatedMeal, err := mealRepo.GetByID(meal.ID)
		require.NoError(t, err)
		assert.Equal(t, 330.0, updatedMeal.TotalCalories)
		assert.Equal(t, 20.0, updatedMeal.TotalProtein)

		// Delete meal
		err = mealRepo.Delete(meal.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = mealRepo.GetByID(meal.ID)
		assert.Error(t, err)
	})
}

func TestMealWithCustomFood(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "custom_food_test@example.com")

	// Initialize repositories
	foodRepo := postgres.NewFoodRepository(testDB.DB)
	mealRepo := postgres.NewMealRepository(testDB.DB)

	t.Run("Create custom food and log meal", func(t *testing.T) {
		// Create custom food
		fiber := 5.0
		customFood := &domain.Food{
			Name:          "My Custom Recipe",
			ServingSize:   100,
			ServingUnit:   "g",
			Calories:      250.0,
			Protein:       15.0,
			Carbohydrates: 30.0,
			Fat:           8.0,
			Fiber:         &fiber,
			IsVerified:    false,
		}

		err := foodRepo.Create(customFood)
		require.NoError(t, err)
		assert.NotEqual(t, "", customFood.ID.String())

		// Verify fiber constraint (fiber < carbs)
		assert.Less(t, *customFood.Fiber, customFood.Carbohydrates)

		// Create meal with custom food
		meal := &domain.Meal{
			UserID:             user.ID,
			Name:               "Dinner",
			MealType:           "dinner",
			ConsumedAt:         time.Now(),
			TotalCalories:      0,
			TotalProtein:       0,
			TotalCarbohydrates: 0,
			TotalFat:           0,
		}

		err = mealRepo.Create(meal)
		require.NoError(t, err)

		// Add custom food to meal
		foodItem := &domain.MealFoodItem{
			MealID:        meal.ID,
			FoodID:        customFood.ID,
			Quantity:      150, // 150g
			Unit:          "g",
			Calories:      375.0, // 250 * 1.5
			Protein:       22.5,  // 15 * 1.5
			Carbohydrates: 45.0,  // 30 * 1.5
			Fat:           12.0,  // 8 * 1.5
		}

		err = testDB.DB.Create(foodItem).Error
		require.NoError(t, err)

		// Update meal totals
		meal.TotalCalories = 375.0
		meal.TotalProtein = 22.5
		meal.TotalCarbohydrates = 45.0
		meal.TotalFat = 12.0
		err = mealRepo.Update(meal)
		require.NoError(t, err)

		// Verify totals
		retrievedMeal, err := mealRepo.GetByID(meal.ID)
		require.NoError(t, err)
		assert.Equal(t, 375.0, retrievedMeal.TotalCalories)
		assert.Equal(t, 22.5, retrievedMeal.TotalProtein)
		assert.Equal(t, 45.0, retrievedMeal.TotalCarbohydrates)
		assert.Equal(t, 12.0, retrievedMeal.TotalFat)
	})
}

func TestGetMealsInRange(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "range_test@example.com")

	// Initialize repository
	mealRepo := postgres.NewMealRepository(testDB.DB)

	t.Run("Get meals within date range", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		twoDaysAgo := now.Add(-48 * time.Hour)

		// Create meals on different days
		meals := []*domain.Meal{
			{
				UserID:             user.ID,
				Name:               "Breakfast Today",
				MealType:           "breakfast",
				ConsumedAt:         now,
				TotalCalories:      300,
				TotalProtein:       15,
				TotalCarbohydrates: 40,
				TotalFat:           10,
			},
			{
				UserID:             user.ID,
				Name:               "Lunch Yesterday",
				MealType:           "lunch",
				ConsumedAt:         yesterday,
				TotalCalories:      500,
				TotalProtein:       30,
				TotalCarbohydrates: 60,
				TotalFat:           15,
			},
			{
				UserID:             user.ID,
				Name:               "Dinner Two Days Ago",
				MealType:           "dinner",
				ConsumedAt:         twoDaysAgo,
				TotalCalories:      600,
				TotalProtein:       40,
				TotalCarbohydrates: 70,
				TotalFat:           20,
			},
		}

		for _, meal := range meals {
			err := mealRepo.Create(meal)
			require.NoError(t, err)
		}

		// Get meals from yesterday onwards
		startDate := yesterday.Add(-1 * time.Hour)
		endDate := now.Add(1 * time.Hour)
		retrievedMeals, err := mealRepo.GetByUserIDAndDateRange(user.ID, startDate, endDate)
		require.NoError(t, err)
		assert.Len(t, retrievedMeals, 2) // Should get today's and yesterday's meals
	})
}
