package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError represents a validation error response
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateRequest creates a middleware that validates request body against a struct
// Usage: router.POST("/api/users", middleware.ValidateRequest(&UserRequest{}), handler)
func ValidateRequest(target interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON to target struct
		if err := c.ShouldBindJSON(&target); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Validate struct
		if err := validate.Struct(target); err != nil {
			validationErrors := make([]ValidationError, 0)

			if validatorErrs, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validatorErrs {
					validationErrors = append(validationErrors, ValidationError{
						Field:   e.Field(),
						Message: getValidationMessage(e),
					})
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Validation failed",
				"validation_errors": validationErrors,
			})
			c.Abort()
			return
		}

		// Store validated object in context
		c.Set("validatedRequest", target)
		c.Next()
	}
}

// getValidationMessage returns a human-readable validation error message
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum: " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum: " + err.Param() + ")"
	case "gte":
		return "Value must be greater than or equal to " + err.Param()
	case "lte":
		return "Value must be less than or equal to " + err.Param()
	case "gt":
		return "Value must be greater than " + err.Param()
	case "lt":
		return "Value must be less than " + err.Param()
	case "len":
		return "Length must be " + err.Param()
	case "oneof":
		return "Value must be one of: " + err.Param()
	case "url":
		return "Invalid URL format"
	case "uuid":
		return "Invalid UUID format"
	default:
		return "Validation failed: " + err.Tag()
	}
}

// ValidateStruct validates a struct and returns validation errors
// This can be used in handlers for manual validation
func ValidateStruct(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors := make([]ValidationError, 0)
	if validatorErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validatorErrs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   e.Field(),
				Message: getValidationMessage(e),
			})
		}
	}

	return validationErrors
}
