package repositories

import (
	"motocare-dashboard/backend/models"

	"gorm.io/gorm"
)

type BookingStatusTotal struct {
	Status string `json:"status"`
	Total  int64  `json:"total"`
}

type TopServiceTotal struct {
	ServiceName   string `json:"service_name"`
	TotalBookings int64  `json:"total_bookings"`
}

type DashboardStats struct {
	TotalCategories   int64                `json:"total_categories"`
	TotalServices     int64                `json:"total_services"`
	TotalBookings     int64                `json:"total_bookings"`
	PendingBookings   int64                `json:"pending_bookings"`
	CompletedBookings int64                `json:"completed_bookings"`
	EstimatedRevenue  int64                `json:"estimated_revenue"`
	BookingsByStatus  []BookingStatusTotal `json:"bookings_by_status"`
	TopServices       []TopServiceTotal    `json:"top_services"`
}

type DashboardRepository interface {
	GetStats(userID uint) (*DashboardStats, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(userID uint) (*DashboardStats, error) {
	stats := &DashboardStats{
		BookingsByStatus: []BookingStatusTotal{},
		TopServices:      []TopServiceTotal{},
	}

	if err := r.db.Model(&models.ServiceCategory{}).Count(&stats.TotalCategories).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.Service{}).Count(&stats.TotalServices).Error; err != nil {
		return nil, err
	}

	bookingQuery := scopedBookingQuery(r.db.Model(&models.Booking{}), userID)
	if err := bookingQuery.Count(&stats.TotalBookings).Error; err != nil {
		return nil, err
	}

	pendingQuery := scopedBookingQuery(r.db.Model(&models.Booking{}), userID).Where("status = ?", "pending")
	if err := pendingQuery.Count(&stats.PendingBookings).Error; err != nil {
		return nil, err
	}

	completedQuery := scopedBookingQuery(r.db.Model(&models.Booking{}), userID).Where("status = ?", "completed")
	if err := completedQuery.Count(&stats.CompletedBookings).Error; err != nil {
		return nil, err
	}

	revenueQuery := scopedBookingQuery(r.db.Model(&models.Booking{}), userID).
		Joins("JOIN services ON services.id = bookings.service_id").
		Where("bookings.status = ?", "completed").
		Select("COALESCE(SUM(services.price), 0)")
	if err := revenueQuery.Scan(&stats.EstimatedRevenue).Error; err != nil {
		return nil, err
	}

	statusQuery := scopedBookingQuery(r.db.Model(&models.Booking{}), userID).
		Select("status, COUNT(*) AS total").
		Group("status").
		Order("status asc")
	if err := statusQuery.Scan(&stats.BookingsByStatus).Error; err != nil {
		return nil, err
	}

	topServicesQuery := scopedBookingQuery(r.db.Model(&models.Booking{}), userID).
		Joins("JOIN services ON services.id = bookings.service_id").
		Select("services.name AS service_name, COUNT(bookings.id) AS total_bookings").
		Group("services.name").
		Order("total_bookings desc").
		Limit(5)
	if err := topServicesQuery.Scan(&stats.TopServices).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func scopedBookingQuery(query *gorm.DB, userID uint) *gorm.DB {
	if userID != 0 {
		return query.Where("bookings.user_id = ?", userID)
	}

	return query
}
