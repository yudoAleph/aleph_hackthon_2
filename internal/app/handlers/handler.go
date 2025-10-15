package handlers

import (
	"net/http"
	"strconv"
	"user-service/internal/app/models"
	"user-service/internal/app/service"
	"user-service/internal/logger"
	"user-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Handler contains methods for handling HTTP requests
type Handler struct {
	service   service.Service
	jwtSecret string
}

func NewHandler(service service.Service, jwtSecret string) *Handler {
	return &Handler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogValidationError(c, "Register", map[string]string{
			"request_body": "Invalid JSON format",
		}, map[string]interface{}{
			"validation_error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	// Validate email format
	if !utils.ValidateEmailField(c, req.Email) {
		logger.LogValidationError(c, "Register", map[string]string{
			"email": "Invalid email format",
		}, map[string]interface{}{
			"email": req.Email,
		})
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		logger.LogEndpointError(c, "Register", err, http.StatusBadRequest, map[string]interface{}{
			"email": req.Email,
		})
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Registration failed",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	// Generate JWT token for the newly registered user
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID

	tokenString, err := token.SignedString([]byte(h.jwtSecret)) // Using the JWT secret from handler
	if err != nil {
		logger.Error(err, map[string]interface{}{
			"handler": "Register",
			"email":   req.Email,
		})
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:     0,
			StatusCode: http.StatusInternalServerError,
			Message:    "Token generation failed",
			Data:       gin.H{"error": "Failed to generate access token"},
		})
		return
	}

	// Create response data with user info and token
	responseData := gin.H{
		"id":         user.ID,
		"full_name":  user.FullName,
		"email":      user.Email,
		"phone":      user.Phone,
		"avatar_url": user.AvatarURL,
		"token": gin.H{
			"access_token": tokenString,
		},
	}

	c.JSON(http.StatusCreated, models.Response{
		Status:     1,
		StatusCode: http.StatusCreated,
		Message:    "Registration success",
		Data:       responseData,
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogValidationError(c, "Login", map[string]string{
			"request_body": "Invalid JSON format",
		}, map[string]interface{}{
			"validation_error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       gin.H{},
		})
		return
	}

	// Validate email format
	if !utils.ValidateEmailField(c, req.Email) {
		logger.LogValidationError(c, "Login", map[string]string{
			"email": "Invalid email format",
		}, map[string]interface{}{
			"email": req.Email,
		})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		logger.LogAuthError(c, "Login", err, map[string]interface{}{
			"email": req.Email,
		})
		c.JSON(http.StatusUnauthorized, models.Response{
			Status:     0,
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid email or password",
			Data:       gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:     1,
		StatusCode: http.StatusOK,
		Message:    "Login success",
		Data:       resp,
	})
}

// GetProfile handles getting the logged-in user's profile
func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		logger.LogEndpointError(c, "GetProfile", err, http.StatusNotFound, map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(http.StatusNotFound, models.Response{
			Status:     0,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
			Data:       gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:     1,
		StatusCode: http.StatusOK,
		Message:    "Profile loaded successfully",
		Data:       user,
	})
}

// UpdateProfile handles updating the logged-in user's profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	userID := c.GetUint("user_id")
	user, err := h.service.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Update failed",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:     1,
		StatusCode: http.StatusOK,
		Message:    "Profile updated successfully",
		Data:       user,
	})
}

// ListContacts handles getting the contact list with search and pagination
func (h *Handler) ListContacts(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.ListContactsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query parameters",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	// Calculate offset for pagination
	req.Offset = (req.Page - 1) * req.Limit

	contacts, count, err := h.service.ListContacts(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:     0,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to load contacts",
			Data:       gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:     1,
		StatusCode: http.StatusOK,
		Message:    "Contacts loaded successfully",
		Data: gin.H{
			"count":    count,
			"page":     req.Page,
			"limit":    req.Limit,
			"contacts": contacts,
		},
	})
}

// CreateContact handles creating a new contact
func (h *Handler) CreateContact(c *gin.Context) {
	var req models.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	// Validate optional email format
	if !utils.ValidateContactEmail(c, req.Email) {
		return
	}

	userID := c.GetUint("user_id")
	contact, err := h.service.CreateContact(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to create contact",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Status:     1,
		StatusCode: http.StatusCreated,
		Message:    "Contact created successfully",
		Data:       contact,
	})
}

// GetContact handles getting a contact's details
func (h *Handler) GetContact(c *gin.Context) {
	userID := c.GetUint("user_id")
	contactID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid contact ID",
			Data:       gin.H{},
		})
		return
	}

	contact, err := h.service.GetContact(c.Request.Context(), userID, uint(contactID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Status:     0,
			StatusCode: http.StatusNotFound,
			Message:    "Contact not found",
			Data:       gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:     1,
		StatusCode: http.StatusOK,
		Message:    "Contact detail loaded",
		Data:       contact,
	})
}

// UpdateContact handles updating a contact
func (h *Handler) UpdateContact(c *gin.Context) {
	var req models.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	// Validate optional email format
	if !utils.ValidateContactEmail(c, req.Email) {
		return
	}

	userID := c.GetUint("user_id")
	contactID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid contact ID",
			Data:       gin.H{},
		})
		return
	}

	contact, err := h.service.UpdateContact(c.Request.Context(), userID, uint(contactID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to update contact",
			Data:       gin.H{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:     1,
		StatusCode: http.StatusOK,
		Message:    "Contact updated successfully",
		Data:       contact,
	})
}

// DeleteContact handles deleting a contact
func (h *Handler) DeleteContact(c *gin.Context) {
	userID := c.GetUint("user_id")
	contactID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:     0,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid contact ID",
			Data:       gin.H{},
		})
		return
	}

	err = h.service.DeleteContact(c.Request.Context(), userID, uint(contactID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Status:     0,
			StatusCode: http.StatusNotFound,
			Message:    "Contact not found",
			Data:       gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:     1,
		StatusCode: http.StatusOK,
		Message:    "Contact deleted successfully",
		Data:       gin.H{},
	})
}
