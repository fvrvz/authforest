package db

import (
	"fmt"
	"log"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DB,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.AuthRefreshTokens{}); err != nil {
		log.Fatalf("Migration failed %v", err)
	}
}

func GetDB() *gorm.DB {
	return db
}
