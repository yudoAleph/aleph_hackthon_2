package utils

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// EmailValidationError represents an email validation error
type EmailValidationError struct {
	Field   string
	Message string
}

// ValidateEmail validates email format using regex
func ValidateEmail(email string) bool {
	// RFC 5322 compliant email regex (simplified)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// ValidateEmailWithResponse validates email and returns appropriate JSON response if invalid
func ValidateEmailWithResponse(c *gin.Context, email string, fieldName string) bool {
	if !ValidateEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      0,
			"status_code": http.StatusBadRequest,
			"message":     "Validation failed",
			"data": gin.H{
				"error": fieldName + " must be a valid email address",
			},
		})
		return false
	}
	return true
}

// ValidateEmailField validates a specific email field from request and returns error response if invalid
func ValidateEmailField(c *gin.Context, email string) bool {
	return ValidateEmailWithResponse(c, email, "email")
}

// ValidateOptionalEmailField validates an optional email field (only if provided)
func ValidateOptionalEmailField(c *gin.Context, email *string, fieldName string) bool {
	if email != nil && *email != "" {
		return ValidateEmailWithResponse(c, *email, fieldName)
	}
	return true
}

// ValidateContactEmail validates email in contact requests (optional field)
func ValidateContactEmail(c *gin.Context, email *string) bool {
	return ValidateOptionalEmailField(c, email, "email")
}
