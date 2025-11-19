package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  ports.UserRepository
	jwtSecret string
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo ports.UserRepository, jwtSecret string) ports.AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	// Validate inputs
	if email == "" || password == "" || name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrDuplicateEntry
	}

	// Hash password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: hashedPassword,
		Name:     name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return password hash
	user.Password = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	// Validate inputs
	if email == "" || password == "" {
		return nil, "", domain.ErrInvalidInput
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, "", domain.ErrUnauthorized
		}
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Compare password
	if err := s.ComparePassword(user.Password, password); err != nil {
		return nil, "", domain.ErrUnauthorized
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Don't return password hash
	user.Password = ""
	return user, token, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	userID, err := s.ParseJWT(tokenString)
	if err != nil {
		return "", domain.ErrUnauthorized
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			return "", domain.ErrUnauthorized
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	return user.ID, nil
}

func (s *authService) HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", domain.ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func (s *authService) ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return domain.ErrUnauthorized
	}
	return nil
}

func (s *authService) GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) ParseJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", domain.ErrUnauthorized
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", domain.ErrUnauthorized
	}

	return userID, nil
}
