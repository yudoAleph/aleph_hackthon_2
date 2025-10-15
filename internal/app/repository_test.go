package app

import (
	"context"
	"testing"

	"user-service/internal/app/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRepository_CreateUser(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful user creation", func(t *testing.T) {
		user := &models.User{
			FullName: "John Doe",
			Email:    "john@example.com",
			Phone:    "+1234567890",
			Password: "hashedpassword",
		}

		createdUser, err := repo.CreateUser(ctx, user)

		require.NoError(t, err)
		assert.NotZero(t, createdUser.ID)
		assert.Equal(t, user.FullName, createdUser.FullName)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.Equal(t, user.Phone, createdUser.Phone)
		assert.Equal(t, user.Password, createdUser.Password)
	})
}

func TestRepository_GetUserByEmail(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("successful user retrieval by email", func(t *testing.T) {
		retrievedUser, err := repo.GetUserByEmail(ctx, user.Email)

		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, retrievedUser.ID)
		assert.Equal(t, user.FullName, retrievedUser.FullName)
		assert.Equal(t, user.Email, retrievedUser.Email)
		assert.Equal(t, user.Phone, retrievedUser.Phone)
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRepository_GetUserByID(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("successful user retrieval by ID", func(t *testing.T) {
		retrievedUser, err := repo.GetUserByID(ctx, createdUser.ID)

		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, retrievedUser.ID)
		assert.Equal(t, user.FullName, retrievedUser.FullName)
		assert.Equal(t, user.Email, retrievedUser.Email)
		assert.Equal(t, user.Phone, retrievedUser.Phone)
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := repo.GetUserByID(ctx, 9999)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRepository_UpdateUser(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("successful user update", func(t *testing.T) {
		updates := map[string]interface{}{
			"full_name": "Updated Name",
			"phone":     "+0987654321",
		}

		updatedUser, err := repo.UpdateUser(ctx, createdUser.ID, updates)

		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, updatedUser.ID)
		assert.Equal(t, "Updated Name", updatedUser.FullName)
		assert.Equal(t, "+0987654321", updatedUser.Phone)
		assert.Equal(t, user.Email, updatedUser.Email) // Email should remain unchanged
	})

	t.Run("user not found for update", func(t *testing.T) {
		updates := map[string]interface{}{
			"full_name": "Updated Name",
		}

		_, err := repo.UpdateUser(ctx, 9999, updates)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRepository_ListContacts(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	// Create test contacts
	contact1 := TestContact(createdUser.ID)
	contact1.FullName = "Alice Johnson"
	contact1.Phone = "+1111111111"

	contact2 := TestContact(createdUser.ID)
	contact2.FullName = "Bob Smith"
	contact2.Phone = "+2222222222"

	_, err = repo.CreateContact(ctx, contact1)
	require.NoError(t, err)

	_, err = repo.CreateContact(ctx, contact2)
	require.NoError(t, err)

	t.Run("list all contacts", func(t *testing.T) {
		contacts, total, err := repo.ListContacts(ctx, createdUser.ID, "", 0, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, contacts, 2)
	})

	t.Run("list contacts with search", func(t *testing.T) {
		contacts, total, err := repo.ListContacts(ctx, createdUser.ID, "Alice", 0, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, contacts, 1)
		assert.Equal(t, "Alice Johnson", contacts[0].FullName)
	})

	t.Run("list contacts with pagination", func(t *testing.T) {
		contacts, total, err := repo.ListContacts(ctx, createdUser.ID, "", 0, 1)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, contacts, 1)
	})

	t.Run("list contacts for non-existent user", func(t *testing.T) {
		contacts, total, err := repo.ListContacts(ctx, 9999, "", 0, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, contacts, 0)
	})
}

func TestRepository_CreateContact(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("successful contact creation", func(t *testing.T) {
		contact := TestContact(createdUser.ID)
		contact.FullName = "Jane Doe"
		contact.Phone = "+3333333333"

		createdContact, err := repo.CreateContact(ctx, contact)

		require.NoError(t, err)
		assert.NotZero(t, createdContact.ID)
		assert.Equal(t, createdUser.ID, createdContact.UserID)
		assert.Equal(t, contact.FullName, createdContact.FullName)
		assert.Equal(t, contact.Phone, createdContact.Phone)
	})
}

func TestRepository_GetContact(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user and contact first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	contact := TestContact(createdUser.ID)
	createdContact, err := repo.CreateContact(ctx, contact)
	require.NoError(t, err)

	t.Run("successful contact retrieval", func(t *testing.T) {
		retrievedContact, err := repo.GetContact(ctx, createdUser.ID, createdContact.ID)

		require.NoError(t, err)
		assert.Equal(t, createdContact.ID, retrievedContact.ID)
		assert.Equal(t, createdUser.ID, retrievedContact.UserID)
		assert.Equal(t, contact.FullName, retrievedContact.FullName)
		assert.Equal(t, contact.Phone, retrievedContact.Phone)
	})

	t.Run("contact not found", func(t *testing.T) {
		_, err := repo.GetContact(ctx, createdUser.ID, 9999)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("contact not found for wrong user", func(t *testing.T) {
		_, err := repo.GetContact(ctx, 9999, createdContact.ID)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRepository_CheckContactExists(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user and contact first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	contact := TestContact(createdUser.ID)
	_, err = repo.CreateContact(ctx, contact)
	require.NoError(t, err)

	t.Run("contact exists", func(t *testing.T) {
		exists, err := repo.CheckContactExists(ctx, createdUser.ID, contact.Phone)

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("contact does not exist", func(t *testing.T) {
		exists, err := repo.CheckContactExists(ctx, createdUser.ID, "+9999999999")

		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("contact does not exist for different user", func(t *testing.T) {
		exists, err := repo.CheckContactExists(ctx, 9999, contact.Phone)

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestRepository_UpdateContact(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user and contact first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	contact := TestContact(createdUser.ID)
	createdContact, err := repo.CreateContact(ctx, contact)
	require.NoError(t, err)

	t.Run("successful contact update", func(t *testing.T) {
		updates := map[string]interface{}{
			"full_name": "Updated Contact",
			"phone":     "+4444444444",
		}

		updatedContact, err := repo.UpdateContact(ctx, createdUser.ID, createdContact.ID, updates)

		require.NoError(t, err)
		assert.Equal(t, createdContact.ID, updatedContact.ID)
		assert.Equal(t, "Updated Contact", updatedContact.FullName)
		assert.Equal(t, "+4444444444", updatedContact.Phone)
	})

	t.Run("contact not found for update", func(t *testing.T) {
		updates := map[string]interface{}{
			"full_name": "Updated Contact",
		}

		_, err := repo.UpdateContact(ctx, createdUser.ID, 9999, updates)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("contact not found for wrong user", func(t *testing.T) {
		updates := map[string]interface{}{
			"full_name": "Updated Contact",
		}

		_, err := repo.UpdateContact(ctx, 9999, createdContact.ID, updates)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRepository_DeleteContact(t *testing.T) {
	_, repo, cleanup := SetupTestEnvironment(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user and contact first
	user := TestUser()
	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	contact := TestContact(createdUser.ID)
	createdContact, err := repo.CreateContact(ctx, contact)
	require.NoError(t, err)

	t.Run("successful contact deletion", func(t *testing.T) {
		err := repo.DeleteContact(ctx, createdUser.ID, createdContact.ID)

		assert.NoError(t, err)

		// Verify contact is deleted
		_, err = repo.GetContact(ctx, createdUser.ID, createdContact.ID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("contact not found for deletion", func(t *testing.T) {
		err := repo.DeleteContact(ctx, createdUser.ID, 9999)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("contact not found for wrong user", func(t *testing.T) {
		// Create another contact for the same user
		contact2 := TestContact(createdUser.ID)
		contact2.Phone = "+5555555555"
		createdContact2, err := repo.CreateContact(ctx, contact2)
		require.NoError(t, err)

		err = repo.DeleteContact(ctx, 9999, createdContact2.ID)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
