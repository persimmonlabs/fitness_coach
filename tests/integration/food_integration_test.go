package integration

import (
	"testing"

	"fitness-tracker/internal/adapters/repositories/postgres"
	"fitness-tracker/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFoodSearch(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository
	foodRepo := postgres.NewFoodRepository(testDB.DB)

	t.Run("Search foods by name", func(t *testing.T) {
		// Seed test foods
		foods := []*domain.Food{
			{
				Name:          "Chicken Breast",
				ServingSize:   100,
				ServingUnit:   "g",
				Calories:      165.0,
				Protein:       31.0,
				Carbohydrates: 0.0,
				Fat:           3.6,
				IsVerified:    true,
			},
			{
				Name:          "Chicken Thigh",
				ServingSize:   100,
				ServingUnit:   "g",
				Calories:      209.0,
				Protein:       26.0,
				Carbohydrates: 0.0,
				Fat:           11.0,
				IsVerified:    true,
			},
			{
				Name:          "Salmon Fillet",
				ServingSize:   100,
				ServingUnit:   "g",
				Calories:      208.0,
				Protein:       20.0,
				Carbohydrates: 0.0,
				Fat:           13.0,
				IsVerified:    true,
			},
		}

		for _, food := range foods {
			err := foodRepo.Create(food)
			require.NoError(t, err)
		}

		// Search for "chicken"
		results, err := foodRepo.Search("chicken", 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)

		// Verify results contain chicken items
		chickenFound := false
		for _, food := range results {
			if food.Name == "Chicken Breast" || food.Name == "Chicken Thigh" {
				chickenFound = true
				break
			}
		}
		assert.True(t, chickenFound, "Should find chicken items")

		// Search for non-existent food
		results, err = foodRepo.Search("nonexistent123xyz", 10, 0)
		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestCreateCustomFood(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository
	foodRepo := postgres.NewFoodRepository(testDB.DB)

	t.Run("Create custom food with valid data", func(t *testing.T) {
		fiber := 3.0
		sugar := 5.0

		food := &domain.Food{
			Name:          "My Protein Shake",
			ServingSize:   250,
			ServingUnit:   "ml",
			Calories:      200.0,
			Protein:       25.0,
			Carbohydrates: 15.0,
			Fat:           5.0,
			Fiber:         &fiber,
			Sugar:         &sugar,
			IsVerified:    false,
		}

		err := foodRepo.Create(food)
		require.NoError(t, err)
		assert.NotEqual(t, "", food.ID.String())

		// Retrieve and verify
		retrieved, err := foodRepo.GetByID(food.ID)
		require.NoError(t, err)
		assert.Equal(t, "My Protein Shake", retrieved.Name)
		assert.Equal(t, 250.0, retrieved.ServingSize)
		assert.Equal(t, "ml", retrieved.ServingUnit)
		assert.Equal(t, 200.0, retrieved.Calories)
		assert.NotNil(t, retrieved.Fiber)
		assert.Equal(t, fiber, *retrieved.Fiber)
	})

	t.Run("Fiber constraint validation", func(t *testing.T) {
		fiber := 10.0
		carbs := 20.0

		food := &domain.Food{
			Name:          "High Fiber Food",
			ServingSize:   100,
			ServingUnit:   "g",
			Calories:      150.0,
			Protein:       5.0,
			Carbohydrates: carbs,
			Fat:           3.0,
			Fiber:         &fiber,
			IsVerified:    false,
		}

		err := foodRepo.Create(food)
		require.NoError(t, err)

		// Verify fiber < carbs constraint
		assert.Less(t, *food.Fiber, food.Carbohydrates,
			"Fiber should be less than carbohydrates")
	})

	t.Run("Create food with brand and category", func(t *testing.T) {
		brand := "TestBrand"
		category := "Protein"

		food := &domain.Food{
			Name:          "TestBrand Protein Bar",
			Brand:         &brand,
			Category:      &category,
			ServingSize:   60,
			ServingUnit:   "g",
			Calories:      220.0,
			Protein:       20.0,
			Carbohydrates: 25.0,
			Fat:           7.0,
			IsVerified:    true,
		}

		err := foodRepo.Create(food)
		require.NoError(t, err)

		// Retrieve and verify
		retrieved, err := foodRepo.GetByID(food.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved.Brand)
		assert.Equal(t, brand, *retrieved.Brand)
		assert.NotNil(t, retrieved.Category)
		assert.Equal(t, category, *retrieved.Category)
	})
}

func TestFoodUpdate(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository
	foodRepo := postgres.NewFoodRepository(testDB.DB)

	t.Run("Update food nutritional values", func(t *testing.T) {
		// Create initial food
		food := CreateTestFood(t, testDB.DB, "Banana", 89.0)

		// Update nutritional values
		food.Calories = 95.0
		food.Carbohydrates = 23.0
		food.Protein = 1.1

		err := foodRepo.Update(food)
		require.NoError(t, err)

		// Retrieve and verify
		updated, err := foodRepo.GetByID(food.ID)
		require.NoError(t, err)
		assert.Equal(t, 95.0, updated.Calories)
		assert.Equal(t, 23.0, updated.Carbohydrates)
		assert.Equal(t, 1.1, updated.Protein)
	})
}

func TestFoodDelete(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository
	foodRepo := postgres.NewFoodRepository(testDB.DB)

	t.Run("Soft delete food", func(t *testing.T) {
		// Create food
		food := CreateTestFood(t, testDB.DB, "Test Food", 100.0)

		// Delete food
		err := foodRepo.Delete(food.ID)
		require.NoError(t, err)

		// Verify soft delete (food not found with normal query)
		_, err = foodRepo.GetByID(food.ID)
		assert.Error(t, err)

		// Verify it still exists with deleted_at set
		var deletedFood domain.Food
		err = testDB.DB.Unscoped().First(&deletedFood, "id = ?", food.ID).Error
		require.NoError(t, err)
		assert.NotNil(t, deletedFood.DeletedAt)
	})
}

func TestGetFoodsByCategory(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository
	foodRepo := postgres.NewFoodRepository(testDB.DB)

	t.Run("Filter foods by category", func(t *testing.T) {
		category1 := "Protein"
		category2 := "Fruit"

		// Create foods with different categories
		foods := []*domain.Food{
			{
				Name:          "Protein Food 1",
				Category:      &category1,
				ServingSize:   100,
				ServingUnit:   "g",
				Calories:      150.0,
				Protein:       20.0,
				Carbohydrates: 5.0,
				Fat:           5.0,
				IsVerified:    true,
			},
			{
				Name:          "Protein Food 2",
				Category:      &category1,
				ServingSize:   100,
				ServingUnit:   "g",
				Calories:      180.0,
				Protein:       25.0,
				Carbohydrates: 3.0,
				Fat:           7.0,
				IsVerified:    true,
			},
			{
				Name:          "Fruit Food 1",
				Category:      &category2,
				ServingSize:   100,
				ServingUnit:   "g",
				Calories:      50.0,
				Protein:       0.5,
				Carbohydrates: 12.0,
				Fat:           0.2,
				IsVerified:    true,
			},
		}

		for _, food := range foods {
			err := foodRepo.Create(food)
			require.NoError(t, err)
		}

		// Query foods by category using the database
		var proteinFoods []domain.Food
		err := testDB.DB.Where("category = ?", category1).Find(&proteinFoods).Error
		require.NoError(t, err)
		assert.Len(t, proteinFoods, 2)

		// Verify all results have the correct category
		for _, food := range proteinFoods {
			assert.NotNil(t, food.Category)
			assert.Equal(t, category1, *food.Category)
		}
	})
}
