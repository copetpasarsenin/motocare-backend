package main

import (
	"fmt"
	"log"
	"motocare-dashboard/backend/config"
	_ "motocare-dashboard/backend/docs"
	"motocare-dashboard/backend/routes"
	"motocare-dashboard/backend/seeders"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/gofiber/swagger"
)

// @title MotoCare Dashboard API
// @version 1.0
// @description REST API for MotoCare Dashboard motorcycle service booking and management.
// @description Use the Authorize button with this format: Bearer <token>.
// @termsOfService http://swagger.io/terms/
// @contact.name MotoCare Dashboard
// @contact.email admin@motocare.test
// @license.name MIT
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func main() {
	config.LoadEnv()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run database migration: %v", err)
	}

	if strings.EqualFold(config.GetEnv("RUN_SEEDER", "false"), "true") {
		if err := seeders.Run(db); err != nil {
			log.Fatalf("failed to run database seeder: %v", err)
		}
		log.Println("database seeder completed")
	}

	app := fiber.New(fiber.Config{
		AppName: "MotoCare Dashboard API",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.GetEnv("FRONTEND_URL", "http://localhost:5173"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.HandlerDefault)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "MotoCare Dashboard API is running",
		})
	})

	routes.SetupAuthRoutes(app, db)
	routes.SetupCRUDRoutes(app, db)

	port := config.GetAppPort()
	log.Printf("server running on port %s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
