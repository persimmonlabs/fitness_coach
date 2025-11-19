package integration

import (
	"testing"
	"time"

	"fitness-tracker/internal/adapters/repositories/postgres"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterLogin(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository and service
	userRepo := postgres.NewUserRepository(testDB.DB)
	authService := services.NewAuthService(userRepo, []byte("test_secret_key"))

	t.Run("Register new user", func(t *testing.T) {
		email := "newuser@example.com"
		password := "SecurePassword123!"

		// Register user
		user, err := authService.Register(email, password, "John", "Doe")
		require.NoError(t, err)
		assert.NotEqual(t, "", user.ID.String())
		assert.Equal(t, email, user.Email)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)

		// Verify password is hashed
		assert.NotEqual(t, password, user.PasswordHash)

		// Verify password can be validated
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		assert.NoError(t, err, "Password should match hash")
	})

	t.Run("Register duplicate email fails", func(t *testing.T) {
		email := "duplicate@example.com"
		password := "Password123!"

		// Register first user
		_, err := authService.Register(email, password, "First", "User")
		require.NoError(t, err)

		// Try to register with same email
		_, err = authService.Register(email, "DifferentPass123!", "Second", "User")
		assert.Error(t, err, "Should fail to register duplicate email")
	})

	t.Run("Login with valid credentials", func(t *testing.T) {
		email := "login@example.com"
		password := "LoginPassword123!"

		// Register user
		_, err := authService.Register(email, password, "Login", "Test")
		require.NoError(t, err)

		// Login
		token, user, err := authService.Login(email, password)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, email, user.Email)
	})

	t.Run("Login with invalid password fails", func(t *testing.T) {
		email := "wrongpass@example.com"
		password := "CorrectPassword123!"

		// Register user
		_, err := authService.Register(email, password, "Wrong", "Pass")
		require.NoError(t, err)

		// Try to login with wrong password
		_, _, err = authService.Login(email, "WrongPassword123!")
		assert.Error(t, err, "Should fail with invalid password")
	})

	t.Run("Login with non-existent user fails", func(t *testing.T) {
		_, _, err := authService.Login("nonexistent@example.com", "Password123!")
		assert.Error(t, err, "Should fail with non-existent user")
	})
}

func TestJWTValidation(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "jwt_test@example.com")

	secretKey := []byte("test_secret_key_for_jwt_validation")

	t.Run("Generate and validate JWT token", func(t *testing.T) {
		// Generate token
		claims := jwt.MapClaims{
			"user_id": user.ID.String(),
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
			"iat":     time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secretKey)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Validate token
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		require.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		// Extract claims
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			assert.Equal(t, user.ID.String(), claims["user_id"])
			assert.Equal(t, user.Email, claims["email"])
		}
	})

	t.Run("Expired token is invalid", func(t *testing.T) {
		// Generate expired token
		claims := jwt.MapClaims{
			"user_id": user.ID.String(),
			"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
			"iat":     time.Now().Add(-2 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secretKey)
		require.NoError(t, err)

		// Try to validate expired token
		_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		assert.Error(t, err, "Should fail with expired token")
	})

	t.Run("Token with invalid signature fails", func(t *testing.T) {
		// Generate token with one key
		claims := jwt.MapClaims{
			"user_id": user.ID.String(),
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secretKey)
		require.NoError(t, err)

		// Try to validate with different key
		wrongKey := []byte("wrong_secret_key")
		_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return wrongKey, nil
		})

		assert.Error(t, err, "Should fail with invalid signature")
	})

	t.Run("Malformed token fails", func(t *testing.T) {
		malformedToken := "not.a.valid.jwt.token"

		_, err := jwt.Parse(malformedToken, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		assert.Error(t, err, "Should fail with malformed token")
	})
}

func TestPasswordHashing(t *testing.T) {
	t.Run("Password is properly hashed", func(t *testing.T) {
		password := "MySecurePassword123!"

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		// Verify hash is different from password
		assert.NotEqual(t, password, string(hash))

		// Verify password matches hash
		err = bcrypt.CompareHashAndPassword(hash, []byte(password))
		assert.NoError(t, err)

		// Verify wrong password doesn't match
		err = bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword"))
		assert.Error(t, err)
	})

	t.Run("Same password generates different hashes", func(t *testing.T) {
		password := "SamePassword123!"

		hash1, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		hash2, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		// Hashes should be different (due to salt)
		assert.NotEqual(t, string(hash1), string(hash2))

		// But both should validate the same password
		assert.NoError(t, bcrypt.CompareHashAndPassword(hash1, []byte(password)))
		assert.NoError(t, bcrypt.CompareHashAndPassword(hash2, []byte(password)))
	})
}

func TestUserProfile(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository
	userRepo := postgres.NewUserRepository(testDB.DB)

	t.Run("Update user profile", func(t *testing.T) {
		// Create user
		user := CreateTestUser(t, testDB.DB, "profile@example.com")

		// Update profile
		height := 175.5
		weight := 75.0
		activityLevel := "moderately_active"

		user.HeightCm = &height
		user.WeightKg = &weight
		user.ActivityLevel = &activityLevel

		err := userRepo.Update(user)
		require.NoError(t, err)

		// Retrieve and verify
		updated, err := userRepo.GetByID(user.ID)
		require.NoError(t, err)
		assert.NotNil(t, updated.HeightCm)
		assert.Equal(t, height, *updated.HeightCm)
		assert.NotNil(t, updated.WeightKg)
		assert.Equal(t, weight, *updated.WeightKg)
		assert.NotNil(t, updated.ActivityLevel)
		assert.Equal(t, activityLevel, *updated.ActivityLevel)
	})
}

func TestGetUserByEmail(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Initialize repository
	userRepo := postgres.NewUserRepository(testDB.DB)

	t.Run("Find user by email", func(t *testing.T) {
		email := "findme@example.com"
		user := CreateTestUser(t, testDB.DB, email)

		// Find by email
		found, err := userRepo.GetByEmail(email)
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, email, found.Email)
	})

	t.Run("Non-existent email returns error", func(t *testing.T) {
		_, err := userRepo.GetByEmail("nonexistent@example.com")
		assert.Error(t, err)
	})
}
