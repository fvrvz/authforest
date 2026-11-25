package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fvrvz/authforest/config"
	"github.com/fvrvz/authforest/constants"
	"github.com/fvrvz/authforest/helpers"
	"github.com/fvrvz/authforest/models"
	"github.com/fvrvz/gologger"
	"github.com/lib/pq"
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

	if err := db.AutoMigrate(&models.User{}, &models.AuthRefreshTokens{}, &models.AccessTokenBlacklist{}, &models.OAuthClient{}, &models.AuthorizationCode{}, &models.Role{}, &models.UserRole{}, &models.PasswordResetToken{}); err != nil {
		log.Fatalf("Migration failed %+v", err)
	}

	// Seed default roles
	adminRole := models.Role{Name: "admin", Description: "Full system administrator"}
	db.FirstOrCreate(&adminRole, models.Role{Name: "admin"})
	userRole := models.Role{Name: "user", Description: "Standard user"}
	db.FirstOrCreate(&userRole, models.Role{Name: "user"})

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
		// Assign admin role to default user
		db.Model(&initialUser).Association("Roles").Replace([]models.Role{adminRole})
	}

	// Seed admin panel OAuth client (public client for the SPA)
	redirectURIMap := make(map[string]bool)
	var redirectURIs []string

	addURI := func(uri string) {
		uri = strings.TrimSpace(uri)
		if uri != "" && !redirectURIMap[uri] {
			redirectURIMap[uri] = true
			redirectURIs = append(redirectURIs, uri)
		}
	}

	// Always support local development default
	addURI("http://localhost:5173/callback")

	// Load from runtime config if specified
	if cfg != nil {
		for _, u := range cfg.AdminPanel.RedirectURIs {
			for _, part := range strings.Split(u, ",") {
				addURI(part)
			}
		}
	}

	// Load from environment variable ADMIN_PANEL_REDIRECT_URIS (comma-separated)
	if envURIs := os.Getenv("ADMIN_PANEL_REDIRECT_URIS"); envURIs != "" {
		for _, part := range strings.Split(envURIs, ",") {
			addURI(part)
		}
	}

	// Load from ADMIN_PANEL_URL or FRONTEND_URL if provided
	for _, envKey := range []string{"ADMIN_PANEL_URL", "FRONTEND_URL"} {
		if u := os.Getenv(envKey); u != "" {
			addURI(strings.TrimRight(u, "/") + "/callback")
		}
	}

	adminClient := models.OAuthClient{
		ClientID:     constants.ADMIN_PANEL_CLIENT_ID,
		ClientName:   constants.ADMIN_PANEL_CLIENT_NAME,
		ClientType:   "public",
		RedirectURIs: redirectURIs,
		Scopes:       "openid profile email",
		GrantTypes:   "authorization_code",
	}

	if err := db.FirstOrCreate(&adminClient, models.OAuthClient{ClientID: constants.ADMIN_PANEL_CLIENT_ID}).Error; err != nil {
		gologger.ERROR("Failed to create admin panel OAuth client: %+v", err)
	} else {
		db.Model(&adminClient).Update("redirect_uris", pq.StringArray(redirectURIs))
		gologger.INFO("Admin panel OAuth client registered with client_id: '%s'", constants.ADMIN_PANEL_CLIENT_ID)
	}

}

func GetDB() *gorm.DB {
	return db
}
