package api

import (
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/abhisheksharm-3/shrtn/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the application router
func SetupRouter(cfg *config.Config) *gin.Engine {
	// Set Gin mode based on config
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Configure CORS with more specific settings
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://yourdomain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.Use(RequestLogger())

	// Initialize services
	urlService := service.NewURLService(cfg)
	analyticsService := service.NewAnalyticsService(cfg)

	// Set up handlers
	urlHandler := NewURLHandler(urlService, analyticsService)

	// API routes
	api := r.Group("/api")
	{
		api.POST("/shorten", urlHandler.ShortenURL)
		api.GET("/urls", urlHandler.GetAllURLs)
		api.GET("/:ShortCode", urlHandler.GetURLByShortCode)    // JSON endpoint
		api.GET("/:ShortCode/redirect", urlHandler.RedirectURL) // Redirect endpoint
	}

	// Direct redirect route (short version)
	r.GET("/:ShortCode", urlHandler.RedirectURL)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return r
}
