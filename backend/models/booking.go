package models

import "time"

type Booking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	User         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	ServiceID    uint      `gorm:"not null;index" json:"service_id"`
	Service      Service   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"service,omitempty"`
	CustomerName string    `gorm:"type:varchar(150);not null" json:"customer_name"`
	Phone        string    `gorm:"type:varchar(30);not null" json:"phone"`
	VehicleName  string    `gorm:"type:varchar(150);not null" json:"vehicle_name"`
	VehiclePlate string    `gorm:"type:varchar(30);not null" json:"vehicle_plate"`
	BookingDate  time.Time `gorm:"not null" json:"booking_date"`
	Status       string    `gorm:"type:varchar(30);not null;default:pending;index" json:"status"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
