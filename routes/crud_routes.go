package routes

import (
	"motocare-dashboard/backend/handlers"
	"motocare-dashboard/backend/middlewares"
	"motocare-dashboard/backend/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupCRUDRoutes(app *fiber.App, db *gorm.DB) {
	categoryRepository := repositories.NewCategoryRepository(db)
	serviceRepository := repositories.NewServiceRepository(db)
	bookingRepository := repositories.NewBookingRepository(db)
	dashboardRepository := repositories.NewDashboardRepository(db)

	categoryHandler := handlers.NewCategoryHandler(categoryRepository)
	serviceHandler := handlers.NewServiceHandler(serviceRepository, categoryRepository)
	bookingHandler := handlers.NewBookingHandler(bookingRepository, serviceRepository)
	dashboardHandler := handlers.NewDashboardHandler(dashboardRepository)

	api := app.Group("/api")

	api.Get("/categories", categoryHandler.List)
	api.Get("/categories/:id", categoryHandler.Detail)
	api.Get("/services", serviceHandler.List)
	api.Get("/services/:id", serviceHandler.Detail)

	admin := api.Group("", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"))
	admin.Post("/categories", categoryHandler.Create)
	admin.Put("/categories/:id", categoryHandler.Update)
	admin.Delete("/categories/:id", categoryHandler.Delete)
	admin.Post("/services", serviceHandler.Create)
	admin.Put("/services/:id", serviceHandler.Update)
	admin.Delete("/services/:id", serviceHandler.Delete)
	admin.Put("/bookings/:id", bookingHandler.UpdateStatus)
	admin.Delete("/bookings/:id", bookingHandler.Delete)

	protected := api.Group("", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin", "user"))
	protected.Get("/bookings", bookingHandler.List)
	protected.Get("/bookings/:id", bookingHandler.Detail)
	protected.Post("/bookings", bookingHandler.Create)
	protected.Get("/dashboard/stats", dashboardHandler.Stats)
}
