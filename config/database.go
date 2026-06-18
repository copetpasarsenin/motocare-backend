package config

import (
	"errors"
	"motocare-dashboard/backend/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("SUPABASE_DSN")
	if dsn == "" {
		return nil, errors.New("SUPABASE_DSN environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.ServiceCategory{},
		&models.Service{},
		&models.Booking{},
	)
}
