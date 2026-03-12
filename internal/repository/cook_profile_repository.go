package repository

import (
	"github.com/yourusername/dished-go-backend/internal/models"
	"gorm.io/gorm"
)

type CookProfileRepository interface {
	Create(profile *models.CookProfile) error
	GetByID(id uint) (*models.CookProfile, error)
	GetAll() ([]models.CookProfile, error)
	Update(profile *models.CookProfile) error
	Delete(id uint) error
}

type cookProfileRepository struct {
	db *gorm.DB
}

func NewCookProfileRepository(db *gorm.DB) CookProfileRepository {
	return &cookProfileRepository{db: db}
}

func (r *cookProfileRepository) Create(profile *models.CookProfile) error {
	return r.db.Create(profile).Error
}

func (r *cookProfileRepository) GetByID(id uint) (*models.CookProfile, error) {
	var profile models.CookProfile
	err := r.db.First(&profile, id).Error
	return &profile, err
}

func (r *cookProfileRepository) GetAll() ([]models.CookProfile, error) {
	var profiles []models.CookProfile
	err := r.db.Find(&profiles).Error
	return profiles, err
}

func (r *cookProfileRepository) Update(profile *models.CookProfile) error {
	return r.db.Save(profile).Error
}

func (r *cookProfileRepository) Delete(id uint) error {
	return r.db.Delete(&models.CookProfile{}, id).Error
}
