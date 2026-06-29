package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSetupAuthRoutesDoesNotProtectLaterPublicRoutes(t *testing.T) {
	app := fiber.New()

	SetupAuthRoutes(app, nil)
	app.Get("/api/public/services", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": []string{}})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/public/services?page=1&limit=100", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected public route status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}
