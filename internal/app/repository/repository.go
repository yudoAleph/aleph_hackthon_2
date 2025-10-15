package repository

import (
	"context"
	"user-service/internal/app/models"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) (*models.User, error)

	ListContacts(ctx context.Context, userID uint, query string, offset, limit int) ([]models.Contact, int64, error)
	CreateContact(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	GetContact(ctx context.Context, userID, contactID uint) (*models.Contact, error)
	CheckContactExists(ctx context.Context, userID uint, phone string) (bool, error)
	UpdateContact(ctx context.Context, userID, contactID uint, updates map[string]interface{}) (*models.Contact, error)
	DeleteContact(ctx context.Context, userID, contactID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateUser creates a new user
func (r *repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *repository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (r *repository) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ListContacts retrieves a paginated list of contacts
func (r *repository) ListContacts(ctx context.Context, userID uint, query string, offset, limit int) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Contact{}).Where("user_id = ?", userID)

	if query != "" {
		db = db.Where("full_name LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&contacts).Error; err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

// CreateContact creates a new contact
func (r *repository) CreateContact(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	if err := r.db.WithContext(ctx).Create(contact).Error; err != nil {
		return nil, err
	}
	return contact, nil
}

// GetContact retrieves a contact by ID and user ID
func (r *repository) GetContact(ctx context.Context, userID, contactID uint) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", contactID, userID).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// CheckContactExists checks if a contact with the given phone number exists for the user
func (r *repository) CheckContactExists(ctx context.Context, userID uint, phone string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Contact{}).
		Where("user_id = ? AND phone = ?", userID, phone).
		Count(&count).Error
	return count > 0, err
}

// UpdateContact updates contact information
func (r *repository) UpdateContact(ctx context.Context, userID, contactID uint, updates map[string]interface{}) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", contactID, userID).First(&contact).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&contact).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &contact, nil
}

// DeleteContact deletes a contact
func (r *repository) DeleteContact(ctx context.Context, userID, contactID uint) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", contactID, userID).Delete(&models.Contact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
