package db

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey;" json:"id"` // cross-db safe
	Username    string     `gorm:"unique;not null" json:"username"`
	Password    string     `gorm:"not null" json:"-"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
