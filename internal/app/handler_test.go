package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/app/handlers"
	"user-service/internal/app/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) Login(ctx context.Context, req models.LoginRequest) (map[string]interface{}, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockService) GetUserProfile(ctx context.Context, userID uint) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) UpdateProfile(ctx context.Context, userID uint, req models.UpdateProfileRequest) (*models.User, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) ListContacts(ctx context.Context, userID uint, req *models.ListContactsRequest) ([]models.Contact, int64, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).([]models.Contact), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) CreateContact(ctx context.Context, userID uint, req *models.CreateContactRequest) (*models.Contact, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockService) GetContact(ctx context.Context, userID, contactID uint) (*models.Contact, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockService) UpdateContact(ctx context.Context, userID, contactID uint, req *models.UpdateContactRequest) (*models.Contact, error) {
	args := m.Called(ctx, userID, contactID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockService) DeleteContact(ctx context.Context, userID, contactID uint) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

func setupTestRouter(mockService *MockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create handler with mock service
	handler := handlers.NewHandler(mockService, "test_secret")

	// Setup routes
	api := router.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)

		protected := api.Group("")
		protected.Use(func(c *gin.Context) {
			// Mock middleware - set user_id in context
			c.Set("user_id", uint(1))
			c.Next()
		})
		{
			protected.GET("/profile", handler.GetProfile)
			protected.PUT("/profile", handler.UpdateProfile)

			protected.GET("/contacts", handler.ListContacts)
			protected.POST("/contacts", handler.CreateContact)
			protected.GET("/contacts/:id", handler.GetContact)
			protected.PUT("/contacts/:id", handler.UpdateContact)
			protected.DELETE("/contacts/:id", handler.DeleteContact)
		}
	}

	return router
}

func TestHandler_Register(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful registration", func(t *testing.T) {
		req := models.RegisterRequest{
			FullName: "John Doe",
			Email:    "john@example.com",
			Phone:    "+1234567890",
			Password: "password123",
		}

		expectedUser := &models.User{
			ID:       1,
			FullName: req.FullName,
			Email:    req.Email,
			Phone:    req.Phone,
		}

		mockService.On("Register", mock.Anything, req).Return(expectedUser, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Registration success", response.Message)
		assert.Contains(t, response.Data, "id")
		assert.Contains(t, response.Data, "token")

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request format", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid request format", response.Message)
	})

	t.Run("registration failed", func(t *testing.T) {
		req := models.RegisterRequest{
			FullName: "Jane Doe",
			Email:    "jane@example.com",
			Phone:    "+0987654321",
			Password: "password123",
		}

		mockService.On("Register", mock.Anything, req).Return(nil, assert.AnError).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Registration failed", response.Message)
	})
}

func TestHandler_Login(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful login", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		loginResponse := map[string]interface{}{
			"id":        float64(1),
			"full_name": "John Doe",
			"email":     req.Email,
			"phone":     "+1234567890",
			"token": map[string]interface{}{
				"access_token": "jwt_token_here",
			},
		}

		mockService.On("Login", mock.Anything, req).Return(loginResponse, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Login success", response.Message)
		assert.Equal(t, loginResponse, response.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request format", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid request format", response.Message)
	})

	t.Run("login failed", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		mockService.On("Login", mock.Anything, req).Return(nil, assert.AnError).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid email or password", response.Message)
	})
}

func TestHandler_GetProfile(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful profile retrieval", func(t *testing.T) {
		userID := uint(1)
		expectedUser := &models.User{
			ID:       userID,
			FullName: "John Doe",
			Email:    "john@example.com",
			Phone:    "+1234567890",
		}

		mockService.On("GetUserProfile", mock.Anything, userID).Return(expectedUser, nil).Once()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/profile", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Profile loaded successfully", response.Message)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(expectedUser.ID), data["id"])
		assert.Equal(t, expectedUser.FullName, data["full_name"])
		assert.Equal(t, expectedUser.Email, data["email"])
		assert.Equal(t, expectedUser.Phone, data["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uint(1)

		mockService.On("GetUserProfile", mock.Anything, userID).Return(nil, assert.AnError).Once()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/profile", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "User not found", response.Message)
	})
}

func TestHandler_UpdateProfile(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful profile update", func(t *testing.T) {
		userID := uint(1)
		req := models.UpdateProfileRequest{
			FullName: "Updated Name",
			Phone:    "+0987654321",
		}

		expectedUser := &models.User{
			ID:       userID,
			FullName: req.FullName,
			Email:    "john@example.com",
			Phone:    req.Phone,
		}

		mockService.On("UpdateProfile", mock.Anything, userID, req).Return(expectedUser, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Profile updated successfully", response.Message)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(expectedUser.ID), data["id"])
		assert.Equal(t, expectedUser.FullName, data["full_name"])
		assert.Equal(t, expectedUser.Email, data["email"])
		assert.Equal(t, expectedUser.Phone, data["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request format", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid request format", response.Message)
	})
}

func TestHandler_ListContacts(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful contact listing", func(t *testing.T) {
		userID := uint(1)
		expectedContacts := []models.Contact{
			{ID: 1, FullName: "Alice", Phone: "+1111111111"},
			{ID: 2, FullName: "Bob", Phone: "+2222222222"},
		}
		expectedTotal := int64(2)

		req := &models.ListContactsRequest{
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListContacts", mock.Anything, userID, req).Return(expectedContacts, expectedTotal, nil).Once()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/contacts?page=1&limit=10", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Contacts loaded successfully", response.Message)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, expectedTotal, int64(data["count"].(float64)))
		assert.Equal(t, float64(1), data["page"])
		assert.Equal(t, float64(10), data["limit"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid query parameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/contacts?page=invalid", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid query parameters", response.Message)
	})
}

func TestHandler_CreateContact(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful contact creation", func(t *testing.T) {
		userID := uint(1)
		req := &models.CreateContactRequest{
			FullName: "New Contact",
			Phone:    "+1234567890",
		}

		expectedContact := &models.Contact{
			ID:       1,
			UserID:   userID,
			FullName: req.FullName,
			Phone:    req.Phone,
		}

		mockService.On("CreateContact", mock.Anything, userID, req).Return(expectedContact, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/contacts", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Contact created successfully", response.Message)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(expectedContact.ID), data["id"])
		assert.Equal(t, expectedContact.FullName, data["full_name"])
		assert.Equal(t, expectedContact.Phone, data["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request format", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/contacts", bytes.NewBufferString("invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid request format", response.Message)
	})
}

func TestHandler_GetContact(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful contact retrieval", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(1)

		expectedContact := &models.Contact{
			ID:       contactID,
			UserID:   userID,
			FullName: "Test Contact",
			Phone:    "+1234567890",
		}

		mockService.On("GetContact", mock.Anything, userID, contactID).Return(expectedContact, nil).Once()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/contacts/1", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Contact detail loaded", response.Message)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(expectedContact.ID), data["id"])
		assert.Equal(t, expectedContact.FullName, data["full_name"])
		assert.Equal(t, expectedContact.Phone, data["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid contact ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/contacts/invalid", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid contact ID", response.Message)
	})

	t.Run("contact not found", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(999)

		mockService.On("GetContact", mock.Anything, userID, contactID).Return(nil, assert.AnError).Once()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/contacts/999", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Contact not found", response.Message)
	})
}

func TestHandler_UpdateContact(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful contact update", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(1)
		req := &models.UpdateContactRequest{
			FullName: "Updated Contact",
			Phone:    "+0987654321",
		}

		expectedContact := &models.Contact{
			ID:       contactID,
			UserID:   userID,
			FullName: req.FullName,
			Phone:    req.Phone,
		}

		mockService.On("UpdateContact", mock.Anything, userID, contactID, req).Return(expectedContact, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/api/contacts/1", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Contact updated successfully", response.Message)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(expectedContact.ID), data["id"])
		assert.Equal(t, expectedContact.FullName, data["full_name"])
		assert.Equal(t, expectedContact.Phone, data["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid contact ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/api/contacts/invalid", bytes.NewBufferString(`{"full_name":"test","phone":"123"}`))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid contact ID", response.Message)
	})

	t.Run("invalid request format", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/api/contacts/1", bytes.NewBufferString("{}"))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid request format", response.Message)
	})
}

func TestHandler_DeleteContact(t *testing.T) {
	mockService := new(MockService)
	router := setupTestRouter(mockService)

	t.Run("successful contact deletion", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(1)

		mockService.On("DeleteContact", mock.Anything, userID, contactID).Return(nil).Once()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/api/contacts/1", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Status)
		assert.Equal(t, "Contact deleted successfully", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid contact ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/api/contacts/invalid", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Invalid contact ID", response.Message)
	})

	t.Run("contact not found", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(999)

		mockService.On("DeleteContact", mock.Anything, userID, contactID).Return(assert.AnError).Once()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/api/contacts/999", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response models.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.Status)
		assert.Equal(t, "Contact not found", response.Message)
	})
}
