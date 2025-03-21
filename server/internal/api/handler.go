package api

import (
	"net/http"

	"github.com/abhisheksharm-3/shrtn/internal/service"

	"github.com/abhisheksharm-3/shrtn/internal/model"
	"github.com/gin-gonic/gin"
)

// URLHandler handles URL shortening requests
type URLHandler struct {
	urlService       *service.URLService
	analyticsService *service.AnalyticsService
}

// NewURLHandler creates a new URLHandler
func NewURLHandler(urlService *service.URLService, analyticsService *service.AnalyticsService) *URLHandler {
	return &URLHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
	}
}

// ShortenURL handles the creation of shortened URLs
func (h *URLHandler) ShortenURL(c *gin.Context) {
	var input model.URLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.urlService.Create(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, url)
}

// RedirectURL redirects to the original URL
func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := h.urlService.GetByShortCode(c, shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Record analytics
	go h.analyticsService.RecordClick(c, url.ID, c.Request)

	// Redirect to original URL
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// GetAllURLs returns all shortened URLs
func (h *URLHandler) GetAllURLs(c *gin.Context) {
	urls, err := h.urlService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, urls)
}
