package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"motocare-dashboard/backend/models"
	"motocare-dashboard/backend/repositories"

	"github.com/gofiber/fiber/v2"
)

type stubServiceRepository struct {
	service models.Service
}

func (r stubServiceRepository) List(params repositories.ServiceListParams) ([]models.Service, int64, error) {
	return []models.Service{r.service}, 1, nil
}

func (r stubServiceRepository) FindByID(id uint) (*models.Service, error) {
	return &r.service, nil
}

func (r stubServiceRepository) Create(service *models.Service) error { return nil }
func (r stubServiceRepository) Update(service *models.Service) error { return nil }
func (r stubServiceRepository) Delete(id uint) error                 { return nil }
func (r stubServiceRepository) Exists(id uint) (bool, error)         { return true, nil }

func TestPublicDetailReturnsPublicShape(t *testing.T) {
	app := fiber.New()
	handler := NewServiceHandler(stubServiceRepository{service: models.Service{
		ID:              7,
		CategoryID:      3,
		Category:        models.ServiceCategory{ID: 3, Name: "Tune Up"},
		Name:            "Servis Ringan",
		Description:     "Perawatan dasar",
		Price:           50000,
		DurationMinutes: 45,
		Status:          "active",
	}}, nil)
	app.Get("/api/public/services/:id", handler.PublicDetail)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/public/services/7", nil))
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var body struct {
		Data map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := body.Data["category_id"]; ok {
		t.Fatalf("public detail leaked internal category_id: %#v", body.Data)
	}
	if body.Data["is_active"] != true {
		t.Fatalf("expected public is_active true, got %#v", body.Data["is_active"])
	}
}
