package utils

import "github.com/gofiber/fiber/v2"

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}

func SuccessResponse(c *fiber.Ctx, status int, message string, data any) error {
	response := fiber.Map{
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	return c.Status(status).JSON(response)
}
