package handlers

import (
	"motocare-dashboard/backend/models"
	"motocare-dashboard/backend/repositories"
	"motocare-dashboard/backend/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ServiceHandler struct {
	serviceRepository  repositories.ServiceRepository
	categoryRepository repositories.CategoryRepository
}

type createServiceRequest struct {
	CategoryID      uint   `json:"category_id" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description"`
	Price           int64  `json:"price" validate:"gte=0"`
	DurationMinutes int    `json:"duration_minutes" validate:"gt=0"`
	Status          string `json:"status" validate:"omitempty,oneof=active inactive"`
}

type updateServiceRequest struct {
	CategoryID      uint   `json:"category_id" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description"`
	Price           int64  `json:"price" validate:"gte=0"`
	DurationMinutes int    `json:"duration_minutes" validate:"gt=0"`
	Status          string `json:"status" validate:"required,oneof=active inactive"`
}

type publicServiceCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type publicServiceResponse struct {
	ID              uint                           `json:"id"`
	Name            string                         `json:"name"`
	Description     string                         `json:"description"`
	Price           int64                          `json:"price"`
	DurationMinutes int                            `json:"duration_minutes"`
	Duration        int                            `json:"duration"`
	Category        *publicServiceCategoryResponse `json:"category,omitempty"`
	Image           *string                        `json:"image"`
	IsActive        bool                           `json:"is_active"`
}

func NewServiceHandler(serviceRepository repositories.ServiceRepository, categoryRepository repositories.CategoryRepository) *ServiceHandler {
	return &ServiceHandler{serviceRepository: serviceRepository, categoryRepository: categoryRepository}
}

func toPublicServiceResponse(service models.Service) publicServiceResponse {
	var category *publicServiceCategoryResponse
	if service.Category.ID != 0 {
		category = &publicServiceCategoryResponse{
			ID:   service.Category.ID,
			Name: service.Category.Name,
		}
	}

	return publicServiceResponse{
		ID:              service.ID,
		Name:            service.Name,
		Description:     service.Description,
		Price:           service.Price,
		DurationMinutes: service.DurationMinutes,
		Duration:        service.DurationMinutes,
		Category:        category,
		Image:           nil,
		IsActive:        service.Status == "active",
	}
}

func (h *ServiceHandler) PublicList(c *fiber.Ctx) error {
	page, limit := utils.ParsePagination(c)

	services, total, err := h.serviceRepository.List(repositories.ServiceListParams{
		Page:      page,
		Limit:     limit,
		Status:    "active",
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil layanan publik")
	}

	publicServices := make([]publicServiceResponse, 0, len(services))
	for _, service := range services {
		publicServices = append(publicServices, toPublicServiceResponse(service))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    publicServices,
		"meta":    utils.NewPaginationMeta(page, limit, total),
	})
}

func (h *ServiceHandler) PublicDetail(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	service, err := h.serviceRepository.FindByID(id)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "layanan tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil layanan publik")
	}
	if service.Status != "active" {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "layanan tidak ditemukan")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "layanan berhasil diambil", toPublicServiceResponse(*service))
}

func (h *ServiceHandler) List(c *fiber.Ctx) error {
	page, limit := utils.ParsePagination(c)

	var categoryID uint
	if raw := strings.TrimSpace(c.Query("category_id")); raw != "" {
		parsed, parseErr := strconv.ParseUint(raw, 10, 64)
		if parseErr != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "category_id tidak valid")
		}
		categoryID = uint(parsed)
	}

	status := strings.TrimSpace(strings.ToLower(c.Query("status")))
	if status != "" && !isAllowedValue(status, "active", "inactive") {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "status hanya boleh active atau inactive")
	}

	services, total, err := h.serviceRepository.List(repositories.ServiceListParams{
		Page:       page,
		Limit:      limit,
		Search:     strings.TrimSpace(c.Query("search")),
		CategoryID: categoryID,
		Status:     status,
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil layanan")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    services,
		"meta":    utils.NewPaginationMeta(page, limit, total),
	})
}

func (h *ServiceHandler) Detail(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	service, err := h.serviceRepository.FindByID(id)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "layanan tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil layanan")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "layanan berhasil diambil", service)
}

func (h *ServiceHandler) Create(c *fiber.Ctx) error {
	if isEmptyBody(c) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "body tidak boleh kosong")
	}

	var request createServiceRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	request.Status = strings.TrimSpace(strings.ToLower(request.Status))
	if request.Status == "" {
		request.Status = "active"
	}

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "validasi gagal", "errors": validationErrors})
	}

	categoryExists, err := h.categoryRepository.Exists(request.CategoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal memeriksa kategori")
	}
	if !categoryExists {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "category_id tidak ditemukan")
	}

	service := models.Service{
		CategoryID:      request.CategoryID,
		Name:            request.Name,
		Description:     request.Description,
		Price:           request.Price,
		DurationMinutes: request.DurationMinutes,
		Status:          request.Status,
	}

	if err := h.serviceRepository.Create(&service); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat layanan")
	}

	createdService, err := h.serviceRepository.FindByID(service.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "layanan dibuat tetapi gagal dimuat")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "layanan berhasil dibuat", createdService)
}

func (h *ServiceHandler) Update(c *fiber.Ctx) error {
	if isEmptyBody(c) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "body tidak boleh kosong")
	}

	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var request updateServiceRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	request.Status = strings.TrimSpace(strings.ToLower(request.Status))

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "validasi gagal", "errors": validationErrors})
	}

	categoryExists, err := h.categoryRepository.Exists(request.CategoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal memeriksa kategori")
	}
	if !categoryExists {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "category_id tidak ditemukan")
	}

	service, err := h.serviceRepository.FindByID(id)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "layanan tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil layanan")
	}

	service.CategoryID = request.CategoryID
	service.Name = request.Name
	service.Description = request.Description
	service.Price = request.Price
	service.DurationMinutes = request.DurationMinutes
	service.Status = request.Status

	if err := h.serviceRepository.Update(service); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengubah layanan")
	}

	updatedService, err := h.serviceRepository.FindByID(service.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "layanan diubah tetapi gagal dimuat")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "layanan berhasil diubah", updatedService)
}

func (h *ServiceHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.serviceRepository.Delete(id); err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "layanan tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal menghapus layanan")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "layanan berhasil dihapus", nil)
}
