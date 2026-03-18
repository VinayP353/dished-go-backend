package repository

import (
	"github.com/yourusername/dished-go-backend/internal/models"
	"gorm.io/gorm"
)

type ChefRepository interface {
	Create(chef *models.Chef) error
	GetByID(id uint) (*models.Chef, error)
	GetByUsername(username string) (*models.Chef, error)
	GetByEmail(email string) (*models.Chef, error)
	GetAll() ([]models.Chef, error)
	Update(chef *models.Chef) error
	UpdateColumns(chef *models.Chef, cols map[string]interface{}) error
	Delete(id uint) error
}

type chefRepository struct {
	db *gorm.DB
}

func NewChefRepository(db *gorm.DB) ChefRepository {
	return &chefRepository{db: db}
}

func (r *chefRepository) Create(chef *models.Chef) error {
	return r.db.Create(chef).Error
}

func (r *chefRepository) GetByID(id uint) (*models.Chef, error) {
	var chef models.Chef
	err := r.db.Preload("ChefProfile").First(&chef, id).Error
	return &chef, err
}

func (r *chefRepository) GetByUsername(username string) (*models.Chef, error) {
	var chef models.Chef
	err := r.db.Preload("ChefProfile").Where("username = ?", username).First(&chef).Error
	return &chef, err
}

func (r *chefRepository) GetByEmail(email string) (*models.Chef, error) {
	var chef models.Chef
	err := r.db.Preload("ChefProfile").Where("email = ?", email).First(&chef).Error
	return &chef, err
}

func (r *chefRepository) GetAll() ([]models.Chef, error) {
	var chefs []models.Chef
	err := r.db.Preload("ChefProfile").Find(&chefs).Error
	return chefs, err
}

func (r *chefRepository) Update(chef *models.Chef) error {
	return r.db.Save(chef).Error
}

func (r *chefRepository) UpdateColumns(chef *models.Chef, cols map[string]interface{}) error {
	return r.db.Model(chef).Updates(cols).Error
}

func (r *chefRepository) Delete(id uint) error {
	return r.db.Delete(&models.Chef{}, id).Error
}
