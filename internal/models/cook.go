package models

import (
	"time"

	"gorm.io/gorm"
)

type Cook struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Username         string         `gorm:"uniqueIndex;not null" json:"username"`
	Password         string         `gorm:"not null" json:"-"`
	Email            string         `gorm:"uniqueIndex;not null" json:"email"`
	Status           string         `gorm:"default:'active'" json:"status"`
	CookieExpiration time.Duration  `gorm:"default:86400000000000" json:"cookie_expiration"` // 24 hours in nanoseconds
	LastLogin        *time.Time     `json:"last_login"`
	CookProfileID    *uint          `json:"cook_profile_id"`
	CookProfile      *CookProfile   `gorm:"foreignKey:CookProfileID" json:"cook_profile,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
