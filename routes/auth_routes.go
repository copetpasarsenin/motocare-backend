package routes

import (
	"motocare-dashboard/backend/handlers"
	"motocare-dashboard/backend/middlewares"
	"motocare-dashboard/backend/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"
)

func rateLimitReached(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
		"message": "terlalu banyak permintaan, coba lagi nanti",
	})
}

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	userRepository := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepository)

	authLimiter := limiter.New(limiter.Config{
		Max:          5,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: rateLimitReached,
	})

	app.Post("/register", authLimiter, authHandler.Register)
	app.Post("/login", authLimiter, authHandler.Login)

	passwordLimiter := limiter.New(limiter.Config{
		Max:          10,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: rateLimitReached,
	})

	protected := app.Group("", middlewares.JWTAuth())
	protected.Get("/me", middlewares.RoleAuthorization("admin", "user"), authHandler.Me)
	protected.Put("/change-password", passwordLimiter, middlewares.RoleAuthorization("admin", "user"), authHandler.ChangePassword)
}
