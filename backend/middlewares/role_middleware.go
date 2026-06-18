package middlewares

import (
	"motocare-dashboard/backend/utils"

	"github.com/gofiber/fiber/v2"
)

func RoleAuthorization(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := GetUserClaims(c)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user belum terautentikasi")
		}

		for _, allowedRole := range allowedRoles {
			if claims.Role == allowedRole {
				return c.Next()
			}
		}

		return utils.ErrorResponse(c, fiber.StatusForbidden, "akses ditolak untuk role ini")
	}
}
