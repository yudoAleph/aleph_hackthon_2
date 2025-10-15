package service

import (
	"context"
	"errors"
	"user-service/internal/app/models"
	"user-service/internal/app/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email is already taken")
	ErrContactNotFound    = errors.New("contact not found")
	ErrPhoneExists        = errors.New("phone number already exists for this user")
)

type Service interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req models.LoginRequest) (map[string]interface{}, error)
	GetUserProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint, req models.UpdateProfileRequest) (*models.User, error)

	ListContacts(ctx context.Context, userID uint, req *models.ListContactsRequest) ([]models.Contact, int64, error)
	CreateContact(ctx context.Context, userID uint, req *models.CreateContactRequest) (*models.Contact, error)
	GetContact(ctx context.Context, userID, contactID uint) (*models.Contact, error)
	UpdateContact(ctx context.Context, userID, contactID uint, req *models.UpdateContactRequest) (*models.Contact, error)
	DeleteContact(ctx context.Context, userID, contactID uint) error
}

type service struct {
	repo      repository.Repository
	jwtSecret string
}

func NewService(repo repository.Repository, jwtSecret string) Service {
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *service) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailTaken
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *service) GetUserProfile(ctx context.Context, userID uint) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID uint, req models.UpdateProfileRequest) (*models.User, error) {
	updates := make(map[string]interface{})
	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	return s.repo.UpdateUser(ctx, userID, updates)
}

func (s *service) ListContacts(ctx context.Context, userID uint, req *models.ListContactsRequest) ([]models.Contact, int64, error) {
	req.Offset = (req.Page - 1) * req.Limit
	return s.repo.ListContacts(ctx, userID, req.Query, req.Offset, req.Limit)
}

func (s *service) CreateContact(ctx context.Context, userID uint, req *models.CreateContactRequest) (*models.Contact, error) {
	// Check if phone number already exists
	exists, err := s.repo.CheckContactExists(ctx, userID, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPhoneExists
	}

	contact := &models.Contact{
		UserID:   userID,
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	return s.repo.CreateContact(ctx, contact)
}

func (s *service) GetContact(ctx context.Context, userID, contactID uint) (*models.Contact, error) {
	contact, err := s.repo.GetContact(ctx, userID, contactID)
	if err != nil {
		return nil, ErrContactNotFound
	}
	return contact, nil
}

func (s *service) UpdateContact(ctx context.Context, userID, contactID uint, req *models.UpdateContactRequest) (*models.Contact, error) {
	// Check if contact exists
	existing, err := s.repo.GetContact(ctx, userID, contactID)
	if err != nil {
		return nil, ErrContactNotFound
	}

	// Check if new phone number conflicts with existing contacts (excluding current contact)
	if existing.Phone != req.Phone {
		exists, err := s.repo.CheckContactExists(ctx, userID, req.Phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrPhoneExists
		}
	}

	updates := map[string]interface{}{
		"full_name": req.FullName,
		"phone":     req.Phone,
		"email":     req.Email,
	}

	return s.repo.UpdateContact(ctx, userID, contactID, updates)
}

func (s *service) DeleteContact(ctx context.Context, userID, contactID uint) error {
	err := s.repo.DeleteContact(ctx, userID, contactID)
	if err != nil {
		return ErrContactNotFound
	}
	return nil
}

// Login authenticates a user and returns a JWT token
func (s *service) Login(ctx context.Context, req models.LoginRequest) (map[string]interface{}, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         user.ID,
		"full_name":  user.FullName,
		"email":      user.Email,
		"phone":      user.Phone,
		"avatar_url": user.AvatarURL,
		"token": models.TokenResponse{
			AccessToken: tokenString,
		},
	}, nil
}
