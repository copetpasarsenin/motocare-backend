package routes

import (
	"motocare-dashboard/backend/handlers"
	"motocare-dashboard/backend/middlewares"
	"motocare-dashboard/backend/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	userRepository := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepository)

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	protected := app.Group("", middlewares.JWTAuth())
	protected.Get("/me", middlewares.RoleAuthorization("admin", "user"), authHandler.Me)
	protected.Put("/change-password", middlewares.RoleAuthorization("admin", "user"), authHandler.ChangePassword)
}
