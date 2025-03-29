package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/abhisheksharm-3/shrtn/internal/model"
	"github.com/abhisheksharm-3/shrtn/internal/service"
	"github.com/gin-gonic/gin"
)

// HTTPError standardizes error responses
type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// URLHandler handles URL shortening requests
type URLHandler struct {
	urlService       *service.URLService
	analyticsService *service.AnalyticsService
	logger           *log.Logger
}

// NewURLHandler creates a new URLHandler
func NewURLHandler(urlService *service.URLService, analyticsService *service.AnalyticsService) *URLHandler {
	return &URLHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
		logger:           log.New(gin.DefaultWriter, "[URL-Handler] ", log.LstdFlags),
	}
}

// ShortenURL handles the creation of shortened URLs
// @Summary Create a shortened URL
// @Description Creates a new shortened URL from the original URL
// @Accept json
// @Produce json
// @Param input body model.URLInput true "URL information"
// @Success 201 {object} model.URL
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/shorten [post]
func (h *URLHandler) ShortenURL(c *gin.Context) {
	h.logger.Println("Processing URL shortening request")

	var input model.URLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_input",
			Message: "Invalid input format: " + err.Error(),
		})
		return
	}

	// Validate URL
	if input.OriginalURL == "" {
		h.handleError(c, HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "missing_url",
			Message: "Original URL cannot be empty",
		})
		return
	}

	// Ensure URL has protocol
	if !strings.HasPrefix(input.OriginalURL, "http://") && !strings.HasPrefix(input.OriginalURL, "https://") {
		input.OriginalURL = "https://" + input.OriginalURL
	}

	// Call the service
	url, err := h.urlService.Create(c, input)
	if err != nil {
		h.logger.Printf("Error creating shortened URL: %v", err)
		h.handleError(c, HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    "creation_failed",
			Message: "Failed to create shortened URL: " + err.Error(),
		})
		return
	}

	h.logger.Printf("Created shortened URL: %s for %s", url.ShortCode, url.OriginalURL)
	c.JSON(http.StatusCreated, url)
}

// GetURLByShortCode returns a URL by its short code without redirecting
// @Summary Get URL by short code
// @Description Retrieves URL information by its short code
// @Produce json
// @Param ShortCode path string true "Short code"
// @Success 200 {object} model.URL
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/{ShortCode} [get]
func (h *URLHandler) GetURLByShortCode(c *gin.Context) {
	shortCode := c.Param("ShortCode")
	if shortCode == "" {
		h.handleError(c, HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "missing_code",
			Message: "Short code cannot be empty",
		})
		return
	}

	url, err := h.urlService.GetByShortCode(c, shortCode)
	if err != nil {
		h.handleError(c, HTTPError{
			Status:  http.StatusNotFound,
			Code:    "not_found",
			Message: "URL not found",
		})
		return
	}

	c.JSON(http.StatusOK, url)
}

// RedirectURL redirects to the original URL
// @Summary Redirect to original URL
// @Description Redirects to the original URL associated with the short code
// @Produce html
// @Param ShortCode path string true "Short code"
// @Success 301 {string} string "Moved Permanently"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /{ShortCode} [get]
func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("ShortCode")
	h.logger.Printf("Processing redirect for: %s", shortCode)

	if shortCode == "" {
		h.handleError(c, HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "missing_code",
			Message: "Short code cannot be empty",
		})
		return
	}

	url, err := h.urlService.GetByShortCode(c, shortCode)
	if err != nil {
		h.logger.Printf("URL not found for code: %s, error: %v", shortCode, err)
		h.handleError(c, HTTPError{
			Status:  http.StatusNotFound,
			Code:    "not_found",
			Message: "URL not found",
		})
		return
	}

	// Ensure URL has protocol
	originalURL := url.OriginalURL
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	// Record analytics asynchronously
	go func() {
		if err := h.analyticsService.RecordClick(c.Request.Context(), url.ID, c.Request); err != nil {
			h.logger.Printf("Failed to record analytics: %v", err)
		}
	}()

	// Update click count asynchronously
	go func() {
		if err := h.urlService.IncrementClicks(c.Request.Context(), url.ID, url.Clicks); err != nil {
			h.logger.Printf("Failed to update click count: %v", err)
		}
	}()

	h.logger.Printf("Redirecting to: %s", originalURL)

	// Set headers to prevent caching
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 301 is permanent redirect, 302 is temporary
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// GetAllURLs returns all shortened URLs
// @Summary Get all URLs
// @Description Retrieves all shortened URLs
// @Produce json
// @Success 200 {array} model.URL
// @Failure 500 {object} HTTPError
// @Router /api/urls [get]
func (h *URLHandler) GetAllURLs(c *gin.Context) {
	h.logger.Println("Retrieving all URLs")

	urls, err := h.urlService.GetAll(c)
	if err != nil {
		h.logger.Printf("Error retrieving URLs: %v", err)
		h.handleError(c, HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    "retrieval_failed",
			Message: fmt.Sprintf("Failed to retrieve URLs: %v", err),
		})
		return
	}

	h.logger.Printf("Retrieved %d URLs", len(urls))
	c.JSON(http.StatusOK, urls)
}

// handleError sends a standardized error response
func (h *URLHandler) handleError(c *gin.Context, err HTTPError) {
	c.JSON(err.Status, gin.H{
		"error": err.Message,
		"code":  err.Code,
	})
}
