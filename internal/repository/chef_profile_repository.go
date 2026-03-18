package repository

import (
	"github.com/yourusername/dished-go-backend/internal/models"
	"gorm.io/gorm"
)

type ChefProfileRepository interface {
	Create(profile *models.ChefProfile) error
	GetByID(id uint) (*models.ChefProfile, error)
	GetAll() ([]models.ChefProfile, error)
	Update(profile *models.ChefProfile) error
	Delete(id uint) error
}

type chefProfileRepository struct {
	db *gorm.DB
}

func NewChefProfileRepository(db *gorm.DB) ChefProfileRepository {
	return &chefProfileRepository{db: db}
}

func (r *chefProfileRepository) Create(profile *models.ChefProfile) error {
	return r.db.Create(profile).Error
}

func (r *chefProfileRepository) GetByID(id uint) (*models.ChefProfile, error) {
	var profile models.ChefProfile
	err := r.db.First(&profile, id).Error
	return &profile, err
}

func (r *chefProfileRepository) GetAll() ([]models.ChefProfile, error) {
	var profiles []models.ChefProfile
	err := r.db.Find(&profiles).Error
	return profiles, err
}

func (r *chefProfileRepository) Update(profile *models.ChefProfile) error {
	return r.db.Save(profile).Error
}

func (r *chefProfileRepository) Delete(id uint) error {
	return r.db.Delete(&models.ChefProfile{}, id).Error
}
