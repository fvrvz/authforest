package server

import (
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fvrvz/authforest/config"
	"github.com/fvrvz/authforest/controllers"
	"github.com/fvrvz/authforest/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// Load HTML templates for the OIDC login/consent pages
	router.LoadHTMLGlob("templates/*")

	setupCors(router)

	// Health check (no auth, no /api/v1 prefix — probed by Docker/orchestrators)
	controllers.SetupHealthRoutes(router)

	// OIDC/OAuth2 routes (at root level per spec)
	controllers.SetupOIDCRoutes(router)

	v1 := router.Group("/api/v1")
	{
		controllers.SetupAuthPublicRoutes(v1)
		controllers.SetupUserPublicRoutes(v1)

		privateGroup := v1.Group("/", middlewares.AuthMiddleware())
		controllers.SetupUserPrivateRoutes(privateGroup)
		controllers.SetupAuthPrivateRoutes(privateGroup)
		controllers.SetupOIDCPrivateRoutes(privateGroup)
	}

	return router
}

func setupCors(router *gin.Engine) {
	corsConfig := config.GetConfig().Server.CORS
	corsDefaults := cors.DefaultConfig()

	originMap := make(map[string]bool)
	var origins []string

	addOrigin := func(origin string) {
		origin = strings.TrimSpace(origin)
		if origin != "" && !originMap[origin] {
			originMap[origin] = true
			origins = append(origins, origin)
		}
	}

	for _, o := range corsConfig.AllowOrigins {
		for _, part := range strings.Split(o, ",") {
			addOrigin(part)
		}
	}

	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		for _, part := range strings.Split(envOrigins, ",") {
			addOrigin(part)
		}
	}

	for _, envKey := range []string{"ADMIN_PANEL_URL", "FRONTEND_URL"} {
		if rawURL := os.Getenv(envKey); rawURL != "" {
			if parsed, err := url.Parse(rawURL); err == nil && parsed.Scheme != "" && parsed.Host != "" {
				addOrigin(parsed.Scheme + "://" + parsed.Host)
			}
		}
	}

	if len(origins) > 0 {
		corsDefaults.AllowOrigins = origins
	}
	if len(corsConfig.AllowMethods) > 0 {
		corsDefaults.AllowMethods = corsConfig.AllowMethods
	}
	if len(corsConfig.AllowHeaders) > 0 {
		corsDefaults.AllowHeaders = corsConfig.AllowHeaders
	}
	if len(corsConfig.ExposeHeaders) > 0 {
		corsDefaults.ExposeHeaders = corsConfig.ExposeHeaders
	}
	if corsConfig.MaxAge > 0 {
		corsDefaults.MaxAge = time.Duration(corsConfig.MaxAge) * time.Hour
	}

	router.Use(cors.New(corsDefaults))
}
