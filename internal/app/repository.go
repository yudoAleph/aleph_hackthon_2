package app

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	Get(ctx context.Context, ID string) (User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Get(ctx context.Context, ID string) (User, error) {
	u := User{}
	err := r.db.WithContext(ctx).Find(&u, "id = ?", ID).Error

	return u, err
}
