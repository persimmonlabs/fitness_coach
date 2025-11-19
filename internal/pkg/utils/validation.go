package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// EmailRegex is a simple email validation regex
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// UsernameRegex validates usernames (alphanumeric, underscore, hyphen, 3-30 chars)
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]{3,30}$`)
)

// ValidationResult holds validation results
type ValidationResult struct {
	Valid  bool
	Errors []string
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:  true,
		Errors: make([]string, 0),
	}
}

// AddError adds an error to the validation result
func (vr *ValidationResult) AddError(err string) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, err)
}

// IsValidEmail validates an email address
func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return EmailRegex.MatchString(email)
}

// IsValidUsername validates a username
func IsValidUsername(username string) bool {
	return UsernameRegex.MatchString(username)
}

// ValidatePassword validates a password with the following requirements:
// - At least 8 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one digit
// - At least one special character
func ValidatePassword(password string) *ValidationResult {
	result := NewValidationResult()

	if len(password) < 8 {
		result.AddError("Password must be at least 8 characters long")
	}

	if len(password) > 128 {
		result.AddError("Password must not exceed 128 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		result.AddError("Password must contain at least one uppercase letter")
	}
	if !hasLower {
		result.AddError("Password must contain at least one lowercase letter")
	}
	if !hasDigit {
		result.AddError("Password must contain at least one digit")
	}
	if !hasSpecial {
		result.AddError("Password must contain at least one special character")
	}

	return result
}

// ValidatePasswordSimple performs basic password validation (minimum length only)
func ValidatePasswordSimple(password string, minLength int) bool {
	if minLength == 0 {
		minLength = 8
	}
	return len(password) >= minLength && len(password) <= 128
}

// SanitizeString removes leading/trailing whitespace and normalizes internal whitespace
func SanitizeString(s string) string {
	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)

	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	return s
}

// IsValidRange checks if a value is within a specified range
func IsValidRange(value, min, max float64) bool {
	return value >= min && value <= max
}

// IsValidIntRange checks if an integer value is within a specified range
func IsValidIntRange(value, min, max int) bool {
	return value >= min && value <= max
}

// IsEmptyString checks if a string is empty or contains only whitespace
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ValidateStringLength validates string length
func ValidateStringLength(s string, minLength, maxLength int) bool {
	length := len(strings.TrimSpace(s))
	return length >= minLength && length <= maxLength
}

// ValidateRequiredFields checks if all required fields are non-empty
func ValidateRequiredFields(fields map[string]string) []string {
	var errors []string
	for fieldName, fieldValue := range fields {
		if IsEmptyString(fieldValue) {
			errors = append(errors, fieldName+" is required")
		}
	}
	return errors
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(urlStr string) bool {
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=]+$`)
	return urlRegex.MatchString(urlStr)
}

// IsValidPhoneNumber validates a phone number (basic validation)
func IsValidPhoneNumber(phone string) bool {
	// Remove common separators
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// Check if it's all digits and has appropriate length
	phoneRegex := regexp.MustCompile(`^\d{10,15}$`)
	return phoneRegex.MatchString(phone)
}

// ValidateAge validates an age value
func ValidateAge(age int) bool {
	return IsValidIntRange(age, 13, 120)
}

// ValidateWeight validates a weight value in kg
func ValidateWeight(weight float64) bool {
	return IsValidRange(weight, 20.0, 500.0)
}

// ValidateHeight validates a height value in cm
func ValidateHeight(height float64) bool {
	return IsValidRange(height, 50.0, 300.0)
}

// IsAlphanumeric checks if a string contains only alphanumeric characters
func IsAlphanumeric(s string) bool {
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphanumericRegex.MatchString(s)
}

// ContainsOnlyLetters checks if a string contains only letters
func ContainsOnlyLetters(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}

// NormalizeEmail normalizes an email address (lowercase, trimmed)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ValidateEnum checks if a value is in the allowed enum values
func ValidateEnum(value string, allowedValues []string) bool {
	for _, allowed := range allowedValues {
		if value == allowed {
			return true
		}
	}
	return false
}

// ValidateEnumCaseInsensitive checks if a value is in the allowed enum values (case-insensitive)
func ValidateEnumCaseInsensitive(value string, allowedValues []string) bool {
	value = strings.ToLower(value)
	for _, allowed := range allowedValues {
		if value == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}
