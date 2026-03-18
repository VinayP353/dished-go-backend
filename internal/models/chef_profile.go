package models

import (
	"time"

	"gorm.io/gorm"
)

type ChefProfile struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	FirstName      string         `gorm:"not null" json:"first_name"`
	LastName       string         `gorm:"not null" json:"last_name"`
	PreferredName  string         `json:"preferred_name"`
	Address        string         `json:"address"`
	ProfilePicture string         `json:"profile_picture"` // Store file path (.jpg, .png)
	Description    string         `gorm:"type:text" json:"description"`
	Signature      string         `json:"signature"` // Store file path or signature data
	Verified       bool           `gorm:"default:false" json:"verified"`
	FHSCertificate string         `json:"fhs_certificate"` // Store file path (image/pdf)
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
