package handlers

import (
	"motocare-dashboard/backend/middlewares"
	"motocare-dashboard/backend/repositories"
	"motocare-dashboard/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	dashboardRepository repositories.DashboardRepository
}

func NewDashboardHandler(dashboardRepository repositories.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{dashboardRepository: dashboardRepository}
}

func (h *DashboardHandler) Stats(c *fiber.Ctx) error {
	claims, ok := middlewares.GetUserClaims(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user belum terautentikasi")
	}

	var userID uint
	if claims.Role == "user" {
		userID = claims.UserID
	}

	stats, err := h.dashboardRepository.GetStats(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil statistik dashboard")
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}
