package app

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username    string     `gorm:"unique" json:"username"`
	Password    string     `gorm:"not null" json:"-"` // also excluded
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}
