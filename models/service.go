package models

import "time"

type Service struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	CategoryID      uint            `gorm:"not null;index" json:"category_id"`
	Category        ServiceCategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`
	Name            string          `gorm:"type:varchar(150);not null" json:"name"`
	Description     string          `gorm:"type:text" json:"description"`
	Price           int64           `gorm:"not null" json:"price"`
	DurationMinutes int             `gorm:"not null" json:"duration_minutes"`
	Status          string          `gorm:"type:varchar(20);not null;default:active;index" json:"status"`
	Bookings        []Booking       `gorm:"foreignKey:ServiceID" json:"bookings,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
