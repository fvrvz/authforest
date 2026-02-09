package db

import (
	"fmt"
	"log"
	"time"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/constants"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/fvrvz/gologger"
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
		log.Fatalf("Failed to connect to DB: %+v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.AuthRefreshTokens{}, &models.AccessTokenBlacklist{}); err != nil {
		log.Fatalf("Migration failed %+v", err)
	}

	hashedPass, err := helpers.HashPassword(constants.DEFAULT_USER_PASS)

	if err != nil {
		gologger.ERROR("Unable to hash password for default user: %+v", err)
	}

	initialUser := models.User{
		Username:  constants.DEFAULT_USERNAME,
		Email:     constants.DEFAULT_USER_EMAIL,
		FirstName: constants.DEFAULT_USER_FNAME,
		LastName:  constants.DEFAULT_USER_LNAME,
		Password:  hashedPass,
		DOB:       time.Now(),
	}

	if err := db.FirstOrCreate(&initialUser, models.User{Username: constants.DEFAULT_USERNAME}).Error; err != nil {
		gologger.ERROR("Failed to create default user: %+v", err)
	} else {
		gologger.INFO("Use default user to login. Username: '%s' with Password: '%s'", constants.DEFAULT_USERNAME, constants.DEFAULT_USER_PASS)
	}

}

func GetDB() *gorm.DB {
	return db
}
