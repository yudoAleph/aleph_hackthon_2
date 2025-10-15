package models

import "time"

// User represents the user model
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName  string    `gorm:"type:varchar(255);not null;index:idx_users_full_name" json:"full_name"`
	Email     string    `gorm:"type:varchar(255);unique;not null;index:idx_users_email" json:"email"`
	Phone     string    `gorm:"type:varchar(20);not null;index:idx_users_phone" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	AvatarURL *string   `gorm:"type:varchar(255)" json:"avatar_url"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_users_created_at" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`

	// Relationships
	Contacts []Contact `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"contacts,omitempty"`
}

// Contact represents the contact model
type Contact struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_contacts_user_id" json:"-"`
	FullName  string    `gorm:"type:varchar(255);not null;index:idx_contacts_full_name" json:"full_name"`
	Phone     string    `gorm:"type:varchar(20);not null;index:idx_contacts_phone" json:"phone"`
	Email     *string   `gorm:"type:varchar(255);index:idx_contacts_email" json:"email"`
	Favorite  bool      `gorm:"default:false;index:idx_contacts_favorite" json:"favorite"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_contacts_created_at" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
