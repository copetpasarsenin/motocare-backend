package repositories

import (
	"motocare-dashboard/backend/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	List() ([]models.ServiceCategory, error)
	FindByID(id uint) (*models.ServiceCategory, error)
	Create(category *models.ServiceCategory) error
	Update(category *models.ServiceCategory) error
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) List() ([]models.ServiceCategory, error) {
	var categories []models.ServiceCategory
	if err := r.db.Order("name asc").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(id uint) (*models.ServiceCategory, error) {
	var category models.ServiceCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Create(category *models.ServiceCategory) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *models.ServiceCategory) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	result := r.db.Delete(&models.ServiceCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *categoryRepository) Exists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.ServiceCategory{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
