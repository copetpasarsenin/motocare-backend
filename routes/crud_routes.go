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

	api.Get("/public/services", serviceHandler.PublicList)
	api.Get("/public/services/:id", serviceHandler.PublicDetail)
	api.Get("/categories", categoryHandler.List)
	api.Get("/categories/:id", categoryHandler.Detail)
	api.Get("/services", serviceHandler.List)
	api.Get("/services/:id", serviceHandler.Detail)

	// USER + ADMIN routes (register first to avoid conflicts)
	api.Get("/bookings", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin", "user"), bookingHandler.List)
	api.Get("/bookings/reserved-slots", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin", "user"), bookingHandler.ReservedSlots)
	api.Get("/bookings/:id", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin", "user"), bookingHandler.Detail)
	api.Post("/bookings", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin", "user"), bookingHandler.Create)

	// ADMIN-only routes
	api.Post("/categories", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), categoryHandler.Create)
	api.Put("/categories/:id", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), categoryHandler.Update)
	api.Delete("/categories/:id", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), categoryHandler.Delete)
	api.Post("/services", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), serviceHandler.Create)
	api.Put("/services/:id", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), serviceHandler.Update)
	api.Delete("/services/:id", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), serviceHandler.Delete)
	api.Put("/bookings/:id", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin", "user"), bookingHandler.UpdateStatus)
	api.Delete("/bookings/:id", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), bookingHandler.Delete)
	api.Get("/dashboard/stats", middlewares.JWTAuth(), middlewares.RoleAuthorization("admin"), dashboardHandler.Stats)
}
