package api

import (
	"fmt"
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
	gin.DefaultWriter.Write([]byte("[DEBUG] Initializing URLHandler\n"))
	return &URLHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
	}
}

// ShortenURL handles the creation of shortened URLs
func (h *URLHandler) ShortenURL(c *gin.Context) {
	requestID := fmt.Sprintf("%p", c.Request)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] ShortenURL: Beginning request processing\n", requestID)))

	// Log headers for debugging CORS or content-type issues
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Request headers: %v\n", requestID, c.Request.Header)))

	// Log request body size
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Request content length: %d\n", requestID, c.Request.ContentLength)))

	var input model.URLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Invalid input format: %s\n", requestID, err.Error())))
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Request body could not be parsed as JSON\n", requestID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log input details after successful parsing
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Parsed input - OriginalURL: %s, CustomCode: %s\n",
		requestID, input.OriginalURL, input.CustomCode)))

	// Validate URL
	if input.OriginalURL == "" {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Empty original URL\n", requestID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Original URL cannot be empty"})
		return
	}

	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Attempting to create shortened URL for: %s\n", requestID, input.OriginalURL)))

	// Call the service
	url, err := h.urlService.Create(c, input)
	if err != nil {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Failed to create shortened URL: %s\n", requestID, err.Error())))

		// Log internal error details for debugging
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Error context: input=%+v\n", requestID, input)))

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Successfully created shortened URL: %s for original URL: %s\n",
		requestID, url.ShortCode, url.OriginalURL)))

	// Log the complete response object (for debugging)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Response: ID=%s, ShortCode=%s, OriginalURL=%s\n",
		requestID, url.ID, url.ShortCode, url.OriginalURL)))

	c.JSON(http.StatusCreated, url)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] ShortenURL: Completed request processing\n", requestID)))
}

// RedirectURL redirects to the original URL
func (h *URLHandler) RedirectURL(c *gin.Context) {
	requestID := fmt.Sprintf("%p", c.Request)
	shortCode := c.Param("ShortCode")

	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] RedirectURL: Beginning request processing for shortCode: %s\n", requestID, shortCode)))

	if shortCode == "" {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Empty short code provided\n", requestID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code cannot be empty"})
		return
	}

	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Looking up original URL for shortCode: %s\n", requestID, shortCode)))

	url, err := h.urlService.GetByShortCode(c, shortCode)
	if err != nil {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Failed to find URL for shortCode %s: %s\n",
			requestID, shortCode, err.Error())))
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Found original URL: %s for shortCode: %s\n",
		requestID, url.OriginalURL, shortCode)))

	// // Record analytics (non-blocking)
	// gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Recording analytics for shortCode: %s\n", requestID, shortCode)))
	// go func() {
	// 	analyticsErr := h.analyticsService.RecordClick(c, url.ID, c.Request)
	// 	if analyticsErr != nil {
	// 		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Failed to record analytics: %s\n",
	// 			requestID, analyticsErr.Error())))
	// 	}
	// }()

	// Redirect to original URL
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Redirecting to: %s\n", requestID, url.OriginalURL)))
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] RedirectURL: Completed request processing\n", requestID)))
}

// GetAllURLs returns all shortened URLs
func (h *URLHandler) GetAllURLs(c *gin.Context) {
	requestID := fmt.Sprintf("%p", c.Request)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] GetAllURLs: Beginning request processing\n", requestID)))

	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Retrieving all URLs\n", requestID)))

	urls, err := h.urlService.GetAll(c)
	if err != nil {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR][%s] Failed to retrieve URLs: %s\n",
			requestID, err.Error())))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] Successfully retrieved %d URLs\n", requestID, len(urls))))

	// Log a sample of the data (first few URLs) if any exist
	if len(urls) > 0 {
		sampleSize := min(3, len(urls))
		for i := 0; i < sampleSize; i++ {
			gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] URL sample %d: ID=%s, ShortCode=%s\n",
				requestID, i, urls[i].ID, urls[i].ShortCode)))
		}
	}

	c.JSON(http.StatusOK, urls)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("[DEBUG][%s] GetAllURLs: Completed request processing\n", requestID)))
}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
