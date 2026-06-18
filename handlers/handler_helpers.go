package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func isEmptyBody(c *fiber.Ctx) bool {
	body := strings.TrimSpace(string(c.Body()))
	return body == "" || body == "{}"
}

func parseIDParam(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "id tidak valid")
	}

	return uint(id), nil
}

func parseUintQuery(c *fiber.Ctx, key string) (uint, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return 0, nil
	}

	parsedValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, key+" tidak valid")
	}

	return uint(parsedValue), nil
}

func parseBookingDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now(), nil
	}

	bookingDate, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return bookingDate, nil
	}

	return time.Parse("2006-01-02", value)
}

func isAllowedValue(value string, allowedValues ...string) bool {
	for _, allowedValue := range allowedValues {
		if value == allowedValue {
			return true
		}
	}

	return false
}
