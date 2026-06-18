package handlers

import (
	"motocare-dashboard/backend/middlewares"
	"motocare-dashboard/backend/models"
	"motocare-dashboard/backend/repositories"
	"motocare-dashboard/backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userRepository repositories.UserRepository
}

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=admin user"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type changePasswordRequest struct {
	UserID          uint   `json:"user_id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

type userResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func NewAuthHandler(userRepository repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepository: userRepository}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request registerRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	request.Role = strings.TrimSpace(strings.ToLower(request.Role))

	if request.Role == "" {
		request.Role = "user"
	}

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validasi gagal",
			"errors":  validationErrors,
		})
	}

	if _, err := h.userRepository.FindByUsername(request.Username); err == nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "username sudah digunakan")
	} else if !repositories.IsRecordNotFound(err) {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal memeriksa username")
	}

	if _, err := h.userRepository.FindByEmail(request.Email); err == nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "email sudah digunakan")
	} else if !repositories.IsRecordNotFound(err) {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal memeriksa email")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengenkripsi password")
	}

	user := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
		Role:     request.Role,
	}

	if err := h.userRepository.Create(&user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat user")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "register berhasil", fiber.Map{
		"user": toUserResponse(user),
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request loginRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validasi gagal",
			"errors":  validationErrors,
		})
	}

	user, err := h.userRepository.FindByEmail(request.Email)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "email atau password salah")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data user")
	}

	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "email atau password salah")
	}

	token, err := utils.GenerateToken(*user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat token")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "login berhasil", fiber.Map{
		"token": token,
		"user":  toUserResponse(*user),
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims, ok := middlewares.GetUserClaims(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user belum terautentikasi")
	}

	user, err := h.userRepository.FindByID(claims.UserID)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data user")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "data user berhasil diambil", fiber.Map{
		"user": toUserResponse(*user),
	})
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	claims, ok := middlewares.GetUserClaims(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user belum terautentikasi")
	}

	var request changePasswordRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "request body tidak valid")
	}

	if validationErrors := utils.ValidateStruct(request); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validasi gagal",
			"errors":  validationErrors,
		})
	}

	targetUserID := claims.UserID
	if claims.Role == "admin" && request.UserID != 0 {
		targetUserID = request.UserID
	}

	if claims.Role != "admin" && request.UserID != 0 && request.UserID != claims.UserID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "user hanya dapat mengubah password sendiri")
	}

	user, err := h.userRepository.FindByID(targetUserID)
	if err != nil {
		if repositories.IsRecordNotFound(err) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user tidak ditemukan")
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data user")
	}

	isChangingOwnPassword := targetUserID == claims.UserID
	if isChangingOwnPassword {
		if strings.TrimSpace(request.CurrentPassword) == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "current_password wajib diisi")
		}

		if !utils.CheckPasswordHash(request.CurrentPassword, user.Password) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "current_password salah")
		}
	}

	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengenkripsi password")
	}

	if err := h.userRepository.UpdatePassword(targetUserID, hashedPassword); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengubah password")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "password berhasil diubah", nil)
}

func toUserResponse(user models.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
