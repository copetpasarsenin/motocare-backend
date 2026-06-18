package handlers

import (
	"motocare-dashboard/backend/middlewares"
	"motocare-dashboard/backend/models"
	"motocare-dashboard/backend/repositories"
	"motocare-dashboard/backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var allowedBookingStatuses = []string{"pending", "confirmed", "in_progress", "completed", "cancelled"}

type BookingHandler struct {
	bookingRepository repositories.BookingRepository
	serviceRepository repositories.ServiceRepository
}

type createBookingRequest struct {
	ServiceID    uint   `json:"service_id" validate:"required"`
	CustomerName string `json:"customer_name" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	VehicleName  string `json:"vehicle_name" validate:"required"`
	VehiclePlate string `json:"vehicle_plate" validate:"required"`
	BookingDate  string `json:"booking_date"`
	Status       string `json:"status" validate:"omitempty,oneof=pending confirmed in_progress completed cancelled"`
	Notes        string `json:"notes"`
}

type updateBookingStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed in_progress completed cancelled"`
}

func NewBookingHandler(bookingRepository repositories.BookingRepository, serviceRepository repositories.ServiceRepository) *BookingHandler {
	return &BookingHandler{bookingRepository: bookingRepository, serviceRepository: serviceRepository}
}

func (h *BookingHandler) List(c *fiber.Ctx) error {
	claims, ok := middlewares.GetUserClaims(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user belum terautentikasi")
	}

	page, limit := utils.ParsePagination(c)
	status := strings.TrimSpace(strings.ToLower(c.Query("status")))
	if status != "" && !isAllowedValue(status, allowedBookingStatuses...) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "status booking tidak valid")
	}

	var userID uint
	if claims.Role == "user" {
		userID = claims.UserID
	}

	bookings, total, err := h.bookingRepository.List(repositories.BookingListParams{
		Page:      page,
		Limit:     limit,
		Search:    strings.TrimSpace(c.Query("search")),
		Status:    status,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		UserID:    userID,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil booking")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    bookings,
		"meta":    utils.NewPaginationMeta(page, limit, total),
	})
}

func (h *BookingHandler) Detail(c *fiber.Ctx) error {
	claims, ok := middlewares.GetUserClaims(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user belum terautentikasi")
	}

	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	booking, err := h.bookingRepository.FindByID(id)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "booking tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil booking")
	}

	if claims.Role == "user" && booking.UserID != claims.UserID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "user hanya dapat melihat booking sendiri")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "booking berhasil diambil", booking)
}

func (h *BookingHandler) Create(c *fiber.Ctx) error {
	claims, ok := middlewares.GetUserClaims(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user belum terautentikasi")
	}

	if isEmptyBody(c) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "body tidak boleh kosong")
	}

	var request createBookingRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.CustomerName = strings.TrimSpace(request.CustomerName)
	request.Phone = strings.TrimSpace(request.Phone)
	request.VehicleName = strings.TrimSpace(request.VehicleName)
	request.VehiclePlate = strings.TrimSpace(strings.ToUpper(request.VehiclePlate))
	request.Status = strings.TrimSpace(strings.ToLower(request.Status))
	request.Notes = strings.TrimSpace(request.Notes)

	if request.Status == "" {
		request.Status = "pending"
	}

	if claims.Role == "user" && request.Status != "pending" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "user hanya dapat membuat booking dengan status pending")
	}

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "validasi gagal", "errors": validationErrors})
	}

	serviceExists, err := h.serviceRepository.Exists(request.ServiceID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal memeriksa layanan")
	}
	if !serviceExists {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "service_id tidak ditemukan")
	}

	bookingDate, err := parseBookingDate(request.BookingDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "booking_date harus format RFC3339 atau YYYY-MM-DD")
	}

	booking := models.Booking{
		UserID:       claims.UserID,
		ServiceID:    request.ServiceID,
		CustomerName: request.CustomerName,
		Phone:        request.Phone,
		VehicleName:  request.VehicleName,
		VehiclePlate: request.VehiclePlate,
		BookingDate:  bookingDate,
		Status:       request.Status,
		Notes:        request.Notes,
	}

	if err := h.bookingRepository.Create(&booking); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat booking")
	}

	createdBooking, err := h.bookingRepository.FindByID(booking.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "booking dibuat tetapi gagal dimuat")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "booking berhasil dibuat", createdBooking)
}

func (h *BookingHandler) UpdateStatus(c *fiber.Ctx) error {
	if isEmptyBody(c) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "body tidak boleh kosong")
	}

	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var request updateBookingStatusRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.Status = strings.TrimSpace(strings.ToLower(request.Status))
	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "validasi gagal", "errors": validationErrors})
	}

	booking, err := h.bookingRepository.UpdateStatus(id, request.Status)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "booking tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengubah status booking")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "status booking berhasil diubah", booking)
}

func (h *BookingHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.bookingRepository.Delete(id); err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "booking tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal menghapus booking")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "booking berhasil dihapus", nil)
}
