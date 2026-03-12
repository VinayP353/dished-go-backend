package repository

import (
	"github.com/yourusername/dished-go-backend/internal/models"
	"gorm.io/gorm"
)

type CookRepository interface {
	Create(cook *models.Cook) error
	GetByID(id uint) (*models.Cook, error)
	GetByUsername(username string) (*models.Cook, error)
	GetByEmail(email string) (*models.Cook, error)
	GetAll() ([]models.Cook, error)
	Update(cook *models.Cook) error
	Delete(id uint) error
}

type cookRepository struct {
	db *gorm.DB
}

func NewCookRepository(db *gorm.DB) CookRepository {
	return &cookRepository{db: db}
}

func (r *cookRepository) Create(cook *models.Cook) error {
	return r.db.Create(cook).Error
}

func (r *cookRepository) GetByID(id uint) (*models.Cook, error) {
	var cook models.Cook
	err := r.db.Preload("CookProfile").First(&cook, id).Error
	return &cook, err
}

func (r *cookRepository) GetByUsername(username string) (*models.Cook, error) {
	var cook models.Cook
	err := r.db.Preload("CookProfile").Where("username = ?", username).First(&cook).Error
	return &cook, err
}

func (r *cookRepository) GetByEmail(email string) (*models.Cook, error) {
	var cook models.Cook
	err := r.db.Preload("CookProfile").Where("email = ?", email).First(&cook).Error
	return &cook, err
}

func (r *cookRepository) GetAll() ([]models.Cook, error) {
	var cooks []models.Cook
	err := r.db.Preload("CookProfile").Find(&cooks).Error
	return cooks, err
}

func (r *cookRepository) Update(cook *models.Cook) error {
	return r.db.Save(cook).Error
}

func (r *cookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cook{}, id).Error
}
