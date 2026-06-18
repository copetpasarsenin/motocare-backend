package repositories

import (
	"fmt"
	"motocare-dashboard/backend/models"
	"strings"

	"gorm.io/gorm"
)

type ServiceListParams struct {
	Page       int
	Limit      int
	Search     string
	CategoryID uint
	Status     string
	SortBy     string
	SortOrder  string
}

type ServiceRepository interface {
	List(params ServiceListParams) ([]models.Service, int64, error)
	FindByID(id uint) (*models.Service, error)
	Create(service *models.Service) error
	Update(service *models.Service) error
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) List(params ServiceListParams) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).Preload("Category")
	query = applyServiceFilters(query, params)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := safeSortBy(params.SortBy, map[string]string{
		"id":               "id",
		"name":             "name",
		"price":            "price",
		"duration_minutes": "duration_minutes",
		"status":           "status",
		"created_at":       "created_at",
	}, "created_at")
	sortOrder := safeSortOrder(params.SortOrder)
	offset := (params.Page - 1) * params.Limit

	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Limit(params.Limit).Offset(offset).Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *serviceRepository) FindByID(id uint) (*models.Service, error) {
	var service models.Service
	if err := r.db.Preload("Category").First(&service, id).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *serviceRepository) Create(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) Update(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Service{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *serviceRepository) Exists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Service{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func applyServiceFilters(query *gorm.DB, params ServiceListParams) *gorm.DB {
	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	if params.CategoryID != 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	return query
}
