package repositories

import (
	"fmt"
	"motocare-dashboard/backend/models"
	"strings"

	"gorm.io/gorm"
)

type BookingListParams struct {
	Page      int
	Limit     int
	Search    string
	Status    string
	SortBy    string
	SortOrder string
	UserID    uint
}

type BookingRepository interface {
	List(params BookingListParams) ([]models.Booking, int64, error)
	FindByID(id uint) (*models.Booking, error)
	Create(booking *models.Booking) error
	UpdateStatus(id uint, status string) (*models.Booking, error)
	Delete(id uint) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) List(params BookingListParams) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).Preload("User").Preload("Service").Preload("Service.Category")
	query = applyBookingFilters(query, params)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := safeSortBy(params.SortBy, map[string]string{
		"id":            "id",
		"customer_name": "customer_name",
		"booking_date":  "booking_date",
		"status":        "status",
		"created_at":    "created_at",
	}, "created_at")
	sortOrder := safeSortOrder(params.SortOrder)
	offset := (params.Page - 1) * params.Limit

	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Limit(params.Limit).Offset(offset).Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (r *bookingRepository) FindByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.Preload("User").Preload("Service").Preload("Service.Category").First(&booking, id).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) UpdateStatus(id uint, status string) (*models.Booking, error) {
	booking, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}

	booking.Status = status
	if err := r.db.Save(booking).Error; err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *bookingRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Booking{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func applyBookingFilters(query *gorm.DB, params BookingListParams) *gorm.DB {
	if params.UserID != 0 {
		query = query.Where("user_id = ?", params.UserID)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(customer_name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(vehicle_name) LIKE ? OR LOWER(vehicle_plate) LIKE ?", search, search, search, search)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	return query
}
