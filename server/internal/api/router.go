// Package api provides HTTP handlers for the URL shortener.
package api

import (
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/abhisheksharm-3/shrtn/internal/middleware"
	"github.com/abhisheksharm-3/shrtn/internal/repository"
	"github.com/abhisheksharm-3/shrtn/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the application router.
func SetupRouter(cfg *config.Config) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.Use(middleware.Security())
	r.Use(middleware.RateLimiter(middleware.RateLimiterConfig{
		RequestsPerMinute: cfg.RateLimitPerMinute,
		BurstSize:         cfg.RateLimitBurst,
		CleanupInterval:   5 * time.Minute,
	}))

	urlRepo := repository.NewAppwriteURLRepository(cfg)
	analyticsRepo := repository.NewAppwriteAnalyticsRepository(cfg)

	urlService := service.NewURLService(urlRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo, "")
	metadataService := service.NewMetadataService()

	urlHandler := NewURLHandler(urlService, analyticsService, metadataService)

	api := r.Group("/api")
	api.Use(middleware.APIKeyAuth(cfg.APIKey))
	{
		api.POST("/shorten", urlHandler.ShortenURL)
		api.GET("/urls", urlHandler.GetAllURLs)
		api.GET("/preview", urlHandler.GetLinkPreview)
		api.GET("/:shortCode", urlHandler.GetURLByShortCode)
		api.DELETE("/:shortCode", urlHandler.DeleteURL)
	}

	r.GET("/:shortCode", urlHandler.RedirectURL)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return r
}
