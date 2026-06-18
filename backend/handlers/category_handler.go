package handlers

import (
	"motocare-dashboard/backend/models"
	"motocare-dashboard/backend/repositories"
	"motocare-dashboard/backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryRepository repositories.CategoryRepository
}

type categoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func NewCategoryHandler(categoryRepository repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categoryRepository: categoryRepository}
}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	categories, err := h.categoryRepository.List()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil kategori")
	}

	total := int64(len(categories))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    categories,
		"meta":    utils.NewPaginationMeta(1, 10, total),
	})
}

func (h *CategoryHandler) Detail(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	category, err := h.categoryRepository.FindByID(id)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kategori tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil kategori")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "kategori berhasil diambil", category)
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	if isEmptyBody(c) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "body tidak boleh kosong")
	}

	var request categoryRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "validasi gagal", "errors": validationErrors})
	}

	category := models.ServiceCategory{Name: request.Name, Description: request.Description}
	if err := h.categoryRepository.Create(&category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat kategori")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "kategori berhasil dibuat", category)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	if isEmptyBody(c) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "body tidak boleh kosong")
	}

	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var request categoryRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "validasi gagal", "errors": validationErrors})
	}

	category, err := h.categoryRepository.FindByID(id)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kategori tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil kategori")
	}

	category.Name = request.Name
	category.Description = request.Description

	if err := h.categoryRepository.Update(category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengubah kategori")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "kategori berhasil diubah", category)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.categoryRepository.Delete(id); err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kategori tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal menghapus kategori")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "kategori berhasil dihapus", nil)
}
