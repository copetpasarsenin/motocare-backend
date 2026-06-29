package config

import (
	"errors"
	"motocare-dashboard/backend/models"
	"os"
	"strconv"
	"time"

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

	if err := configureConnectionPool(db); err != nil {
		return nil, err
	}

	return db, nil
}

func configureConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	maxOpenConns := envInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := envInt("DB_MAX_IDLE_CONNS", 10)
	connMaxLifetime := envDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	connMaxIdleTime := envDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	return sqlDB.Ping()
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.ServiceCategory{},
		&models.Service{},
		&models.Booking{},
	)
}
