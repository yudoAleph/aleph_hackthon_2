package app

import (
	"context"
	"errors"
	"testing"
	"user-service/internal/app/models"
	"user-service/internal/app/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// Test errors (matching service package errors)
var (
	ErrEmailTaken      = errors.New("email is already taken")
	ErrContactNotFound = errors.New("contact not found")
	ErrPhoneExists     = errors.New("phone number already exists for this user")
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) (*models.User, error) {
	args := m.Called(ctx, userID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) ListContacts(ctx context.Context, userID uint, query string, offset, limit int) ([]models.Contact, int64, error) {
	args := m.Called(ctx, userID, query, offset, limit)
	return args.Get(0).([]models.Contact), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) CreateContact(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	args := m.Called(ctx, contact)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockRepository) GetContact(ctx context.Context, userID, contactID uint) (*models.Contact, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockRepository) CheckContactExists(ctx context.Context, userID uint, phone string) (bool, error) {
	args := m.Called(ctx, userID, phone)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) UpdateContact(ctx context.Context, userID, contactID uint, updates map[string]interface{}) (*models.Contact, error) {
	args := m.Called(ctx, userID, contactID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockRepository) DeleteContact(ctx context.Context, userID, contactID uint) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

func TestService_Register(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

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

		// Mock repository calls
		mockRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, nil).Once()
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(expectedUser, nil).Once()

		user, err := service.Register(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.FullName, user.FullName)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.Phone, user.Phone)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already taken", func(t *testing.T) {
		req := models.RegisterRequest{
			FullName: "Jane Doe",
			Email:    "existing@example.com",
			Phone:    "+1234567890",
			Password: "password123",
		}

		existingUser := &models.User{
			ID:    1,
			Email: req.Email,
		}

		mockRepo.On("GetUserByEmail", ctx, req.Email).Return(existingUser, nil).Once()

		user, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, ErrEmailTaken, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_Login(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

	t.Run("successful login", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		// Create a proper bcrypt hash of the password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user := &models.User{
			ID:       1,
			FullName: "John Doe",
			Email:    req.Email,
			Phone:    "+1234567890",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetUserByEmail", ctx, req.Email).Return(user, nil).Once()

		result, err := service.Login(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, user.ID, result["id"])
		assert.Equal(t, user.FullName, result["full_name"])
		assert.Equal(t, user.Email, result["email"])
		assert.Equal(t, user.Phone, result["phone"])
		assert.Contains(t, result, "token")
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, errors.New("user not found")).Once()

		result, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		hashedPassword := "$2a$10$hashedpassword"
		user := &models.User{
			ID:       1,
			Email:    req.Email,
			Password: hashedPassword,
		}

		mockRepo.On("GetUserByEmail", ctx, req.Email).Return(user, nil).Once()

		result, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetUserProfile(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

	t.Run("successful profile retrieval", func(t *testing.T) {
		userID := uint(1)
		expectedUser := &models.User{
			ID:       userID,
			FullName: "John Doe",
			Email:    "john@example.com",
			Phone:    "+1234567890",
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(expectedUser, nil).Once()

		user, err := service.GetUserProfile(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uint(999)

		mockRepo.On("GetUserByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		user, err := service.GetUserProfile(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateProfile(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

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

		mockRepo.On("UpdateUser", ctx, userID, mock.AnythingOfType("map[string]interface {}")).Return(expectedUser, nil).Once()

		user, err := service.UpdateProfile(ctx, userID, req)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_ListContacts(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

	t.Run("successful contact listing", func(t *testing.T) {
		userID := uint(1)
		req := &models.ListContactsRequest{
			Page:  1,
			Limit: 10,
			Query: "test",
		}

		expectedContacts := []models.Contact{
			{ID: 1, FullName: "Test Contact", Phone: "+1234567890"},
		}
		expectedTotal := int64(1)

		mockRepo.On("ListContacts", ctx, userID, req.Query, 0, req.Limit).Return(expectedContacts, expectedTotal, nil).Once()

		contacts, total, err := service.ListContacts(ctx, userID, req)

		require.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CreateContact(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

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

		mockRepo.On("CheckContactExists", ctx, userID, req.Phone).Return(false, nil).Once()
		mockRepo.On("CreateContact", ctx, mock.AnythingOfType("*models.Contact")).Return(expectedContact, nil).Once()

		contact, err := service.CreateContact(ctx, userID, req)

		require.NoError(t, err)
		assert.Equal(t, expectedContact, contact)
		mockRepo.AssertExpectations(t)
	})

	t.Run("phone number already exists", func(t *testing.T) {
		userID := uint(1)
		req := &models.CreateContactRequest{
			FullName: "New Contact",
			Phone:    "+1234567890",
		}

		mockRepo.On("CheckContactExists", ctx, userID, req.Phone).Return(true, nil).Once()

		contact, err := service.CreateContact(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, ErrPhoneExists, err)
		assert.Nil(t, contact)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetContact(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

	t.Run("successful contact retrieval", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(1)

		expectedContact := &models.Contact{
			ID:       contactID,
			UserID:   userID,
			FullName: "Test Contact",
			Phone:    "+1234567890",
		}

		mockRepo.On("GetContact", ctx, userID, contactID).Return(expectedContact, nil).Once()

		contact, err := service.GetContact(ctx, userID, contactID)

		require.NoError(t, err)
		assert.Equal(t, expectedContact, contact)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(999)

		mockRepo.On("GetContact", ctx, userID, contactID).Return(nil, errors.New("contact not found")).Once()

		contact, err := service.GetContact(ctx, userID, contactID)

		assert.Error(t, err)
		assert.Equal(t, ErrContactNotFound, err)
		assert.Nil(t, contact)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateContact(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

	t.Run("successful contact update", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(1)
		req := &models.UpdateContactRequest{
			FullName: "Updated Contact",
			Phone:    "+0987654321",
		}

		existingContact := &models.Contact{
			ID:       contactID,
			UserID:   userID,
			FullName: "Old Contact",
			Phone:    "+1234567890",
		}

		updatedContact := &models.Contact{
			ID:       contactID,
			UserID:   userID,
			FullName: req.FullName,
			Phone:    req.Phone,
		}

		mockRepo.On("GetContact", ctx, userID, contactID).Return(existingContact, nil).Once()
		mockRepo.On("CheckContactExists", ctx, userID, req.Phone).Return(false, nil).Once()
		mockRepo.On("UpdateContact", ctx, userID, contactID, mock.AnythingOfType("map[string]interface {}")).Return(updatedContact, nil).Once()

		contact, err := service.UpdateContact(ctx, userID, contactID, req)

		require.NoError(t, err)
		assert.Equal(t, updatedContact, contact)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(999)
		req := &models.UpdateContactRequest{
			FullName: "Updated Contact",
			Phone:    "+0987654321",
		}

		mockRepo.On("GetContact", ctx, userID, contactID).Return(nil, errors.New("contact not found")).Once()

		contact, err := service.UpdateContact(ctx, userID, contactID, req)

		assert.Error(t, err)
		assert.Equal(t, ErrContactNotFound, err)
		assert.Nil(t, contact)
		mockRepo.AssertExpectations(t)
	})

	t.Run("phone number already exists", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(1)
		req := &models.UpdateContactRequest{
			FullName: "Updated Contact",
			Phone:    "+0987654321",
		}

		existingContact := &models.Contact{
			ID:       contactID,
			UserID:   userID,
			FullName: "Old Contact",
			Phone:    "+1234567890",
		}

		mockRepo.On("GetContact", ctx, userID, contactID).Return(existingContact, nil).Once()
		mockRepo.On("CheckContactExists", ctx, userID, req.Phone).Return(true, nil).Once()

		contact, err := service.UpdateContact(ctx, userID, contactID, req)

		assert.Error(t, err)
		assert.Equal(t, ErrPhoneExists, err)
		assert.Nil(t, contact)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_DeleteContact(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo, "test_secret")
	ctx := context.Background()

	t.Run("successful contact deletion", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(1)

		mockRepo.On("DeleteContact", ctx, userID, contactID).Return(nil).Once()

		err := service.DeleteContact(ctx, userID, contactID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		userID := uint(1)
		contactID := uint(999)

		mockRepo.On("DeleteContact", ctx, userID, contactID).Return(errors.New("contact not found")).Once()

		err := service.DeleteContact(ctx, userID, contactID)

		assert.Error(t, err)
		assert.Equal(t, ErrContactNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
