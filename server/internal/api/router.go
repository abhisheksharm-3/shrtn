package api

import (
	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/abhisheksharm-3/shrtn/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())
	r.Use(cors.Default())
	r.Use(RequestLogger())

	urlService := service.NewURLService(cfg)
	analyticsService := service.NewAnalyticsService(cfg)

	urlHandler := NewURLHandler(urlService, analyticsService)
	api := r.Group("/api")
	{
		api.POST("/shorten", urlHandler.ShortenURL)
		api.GET("/urls", urlHandler.GetAllURLs)
	}

	r.GET("/:shortCode", urlHandler.RedirectURL)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	return r
}
