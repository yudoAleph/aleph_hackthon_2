package models

// TokenResponse represents the token response structure
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// RegisterRequest represents the registration request structure
type RegisterRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    *string `json:"phone,omitempty"`
	Password string  `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the login request structure
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the profile update request structure
type UpdateProfileRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	Phone    *string `json:"phone,omitempty"`
}

// CreateContactRequest represents the create contact request structure
type CreateContactRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	Phone    string  `json:"phone" binding:"required"`
	Email    *string `json:"email"`
}

// UpdateContactRequest represents the update contact request structure
type UpdateContactRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	Phone    string  `json:"phone" binding:"required"`
	Email    *string `json:"email"`
	Favorite bool    `json:"favorite"`
}

// Response represents the standard API response structure
type Response struct {
	Status     int         `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
