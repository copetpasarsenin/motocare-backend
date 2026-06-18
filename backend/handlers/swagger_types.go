package handlers

import "motocare-dashboard/backend/models"

type ErrorResponseDoc struct {
	Message string            `json:"message" example:"validasi gagal"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type MessageResponseDoc struct {
	Message string `json:"message" example:"success"`
}

type MetaResponseDoc struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

type UserDoc struct {
	ID        uint   `json:"id" example:"1"`
	Username  string `json:"username" example:"admin"`
	Email     string `json:"email" example:"admin@motocare.test"`
	Role      string `json:"role" example:"admin"`
	CreatedAt string `json:"created_at,omitempty" example:"2026-06-17T15:04:05+07:00"`
	UpdatedAt string `json:"updated_at,omitempty" example:"2026-06-17T15:04:05+07:00"`
}

type RegisterRequestDoc struct {
	Username string `json:"username" example:"budi"`
	Email    string `json:"email" example:"budi@example.com"`
	Password string `json:"password" example:"password123"`
	Role     string `json:"role,omitempty" example:"user" enums:"admin,user"`
}

type LoginRequestDoc struct {
	Email    string `json:"email" example:"admin@motocare.test"`
	Password string `json:"password" example:"password123"`
}

type ChangePasswordRequestDoc struct {
	UserID          uint   `json:"user_id,omitempty" example:"2"`
	CurrentPassword string `json:"current_password,omitempty" example:"password123"`
	NewPassword     string `json:"new_password" example:"newpassword123"`
}

type UserDataDoc struct {
	User UserDoc `json:"user"`
}

type RegisterResponseDoc struct {
	Message string      `json:"message" example:"register berhasil"`
	Data    UserDataDoc `json:"data"`
}

type LoginDataDoc struct {
	Token string  `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserDoc `json:"user"`
}

type LoginResponseDoc struct {
	Message string       `json:"message" example:"login berhasil"`
	Data    LoginDataDoc `json:"data"`
}

type UserResponseDoc struct {
	Message string      `json:"message" example:"data user berhasil diambil"`
	Data    UserDataDoc `json:"data"`
}

type CategoryRequestDoc struct {
	Name        string `json:"name" example:"Ganti Oli"`
	Description string `json:"description" example:"Layanan penggantian oli mesin dan oli transmisi."`
}

type CategoryResponseDoc struct {
	Message string                 `json:"message" example:"kategori berhasil diambil"`
	Data    models.ServiceCategory `json:"data"`
}

type CategoryListResponseDoc struct {
	Message string                   `json:"message" example:"success"`
	Data    []models.ServiceCategory `json:"data"`
	Meta    MetaResponseDoc          `json:"meta"`
}

type ServiceRequestDoc struct {
	CategoryID      uint   `json:"category_id" example:"1"`
	Name            string `json:"name" example:"Ganti Oli Mesin"`
	Description     string `json:"description" example:"Penggantian oli mesin standar bengkel."`
	Price           int64  `json:"price" example:"65000"`
	DurationMinutes int    `json:"duration_minutes" example:"30"`
	Status          string `json:"status" example:"active" enums:"active,inactive"`
}

type ServiceResponseDoc struct {
	Message string         `json:"message" example:"layanan berhasil diambil"`
	Data    models.Service `json:"data"`
}

type ServiceListResponseDoc struct {
	Message string           `json:"message" example:"success"`
	Data    []models.Service `json:"data"`
	Meta    MetaResponseDoc  `json:"meta"`
}

type BookingRequestDoc struct {
	ServiceID    uint   `json:"service_id" example:"1"`
	CustomerName string `json:"customer_name" example:"Budi Santoso"`
	Phone        string `json:"phone" example:"081234560001"`
	VehicleName  string `json:"vehicle_name" example:"Honda Beat"`
	VehiclePlate string `json:"vehicle_plate" example:"B 1234 MTC"`
	BookingDate  string `json:"booking_date" example:"2026-06-18"`
	Status       string `json:"status,omitempty" example:"pending" enums:"pending,confirmed,in_progress,completed,cancelled"`
	Notes        string `json:"notes" example:"Servis pagi jika tersedia."`
}

type BookingStatusRequestDoc struct {
	Status string `json:"status" example:"confirmed" enums:"pending,confirmed,in_progress,completed,cancelled"`
}

type BookingResponseDoc struct {
	Message string         `json:"message" example:"booking berhasil diambil"`
	Data    models.Booking `json:"data"`
}

type BookingListResponseDoc struct {
	Message string           `json:"message" example:"success"`
	Data    []models.Booking `json:"data"`
	Meta    MetaResponseDoc  `json:"meta"`
}

type DashboardBookingStatusDoc struct {
	Status string `json:"status" example:"pending"`
	Total  int64  `json:"total" example:"3"`
}

type DashboardTopServiceDoc struct {
	ServiceName   string `json:"service_name" example:"Ganti Oli Mesin"`
	TotalBookings int64  `json:"total_bookings" example:"5"`
}

type DashboardStatsDoc struct {
	TotalCategories   int64                       `json:"total_categories" example:"10"`
	TotalServices     int64                       `json:"total_services" example:"10"`
	TotalBookings     int64                       `json:"total_bookings" example:"10"`
	PendingBookings   int64                       `json:"pending_bookings" example:"3"`
	CompletedBookings int64                       `json:"completed_bookings" example:"2"`
	EstimatedRevenue  int64                       `json:"estimated_revenue" example:"500000"`
	BookingsByStatus  []DashboardBookingStatusDoc `json:"bookings_by_status"`
	TopServices       []DashboardTopServiceDoc    `json:"top_services"`
}
