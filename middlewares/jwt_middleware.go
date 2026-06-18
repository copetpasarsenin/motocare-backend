package middlewares

import (
	"motocare-dashboard/backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const UserClaimsKey = "user_claims"

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")
		if authorizationHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "authorization header wajib diisi")
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "format authorization harus Bearer <token>")
		}

		claims, err := utils.ParseToken(strings.TrimSpace(parts[1]))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "token tidak valid atau sudah kedaluwarsa")
		}

		c.Locals(UserClaimsKey, claims)
		return c.Next()
	}
}

func GetUserClaims(c *fiber.Ctx) (*utils.JWTClaims, bool) {
	claims, ok := c.Locals(UserClaimsKey).(*utils.JWTClaims)
	return claims, ok
}
