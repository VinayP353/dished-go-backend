package models

import (
	"time"

	"gorm.io/gorm"
)

type Chef struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Username         string         `gorm:"uniqueIndex;not null" json:"username"`
	Password         string         `gorm:"not null" json:"-"`
	Email            string         `gorm:"uniqueIndex;not null" json:"email"`
	Status           string         `gorm:"default:'active'" json:"status"`
	CookieExpiration int64          `gorm:"default:86400000000000" json:"cookie_expiration"` // nanoseconds
	LastLogin        *time.Time     `json:"last_login"`
	ChefProfileID    *uint          `json:"chef_profile_id"`
	ChefProfile      *ChefProfile   `gorm:"foreignKey:ChefProfileID" json:"chef_profile,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
