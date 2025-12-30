// Package api provides HTTP handlers for the URL shortener.
package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abhisheksharm-3/shrtn/internal/model"
	"github.com/abhisheksharm-3/shrtn/internal/service"
	"github.com/gin-gonic/gin"
)

// URLHandler handles URL shortening HTTP requests.
type URLHandler struct {
	urlService       *service.URLService
	analyticsService *service.AnalyticsService
	metadataService  *service.MetadataService
}

// NewURLHandler creates a new URLHandler.
func NewURLHandler(urlService *service.URLService, analyticsService *service.AnalyticsService, metadataService *service.MetadataService) *URLHandler {
	return &URLHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
		metadataService:  metadataService,
	}
}

// ShortenURL handles POST /api/shorten requests.
func (h *URLHandler) ShortenURL(c *gin.Context) {
	var input model.URLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input format",
			"code":  "invalid_input",
		})
		return
	}

	url, err := h.urlService.Create(c.Request.Context(), input)
	if err != nil {
		status := http.StatusInternalServerError
		code := "creation_failed"

		switch err {
		case service.ErrInvalidURL:
			status = http.StatusBadRequest
			code = "invalid_url"
		case service.ErrURLBlocked:
			status = http.StatusBadRequest
			code = "url_blocked"
		case service.ErrShortCodeExists:
			status = http.StatusConflict
			code = "code_exists"
		case service.ErrShortCodeTooShort, service.ErrShortCodeInvalid:
			status = http.StatusBadRequest
			code = "invalid_code"
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
			"code":  code,
		})
		return
	}

	c.JSON(http.StatusCreated, url)
}

// GetURLByShortCode handles GET /api/:shortCode requests.
func (h *URLHandler) GetURLByShortCode(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "short code is required",
			"code":  "missing_code",
		})
		return
	}

	url, err := h.urlService.GetByShortCode(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
			"code":  "not_found",
		})
		return
	}

	c.JSON(http.StatusOK, url)
}

// RedirectURL handles GET /:shortCode requests.
func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "short code is required",
			"code":  "missing_code",
		})
		return
	}

	url, err := h.urlService.GetByShortCode(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
			"code":  "not_found",
		})
		return
	}

	bgCtx := context.Background()
	go func() {
		_ = h.analyticsService.RecordClick(bgCtx, url.ID, c.Request)
	}()
	go func() {
		_ = h.urlService.IncrementClicks(bgCtx, url.ID, url.Clicks)
	}()

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// GetAllURLs handles GET /api/urls requests.
func (h *URLHandler) GetAllURLs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	response, err := h.urlService.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve URLs",
			"code":  "retrieval_failed",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteURL handles DELETE /api/:shortCode requests.
func (h *URLHandler) DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "short code is required",
			"code":  "missing_code",
		})
		return
	}

	url, err := h.urlService.GetByShortCode(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
			"code":  "not_found",
		})
		return
	}

	if err := h.urlService.Delete(c.Request.Context(), url.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete URL",
			"code":  "deletion_failed",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetLinkPreview handles GET /api/preview requests.
func (h *URLHandler) GetLinkPreview(c *gin.Context) {
	targetURL := c.Query("url")
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "url query parameter is required",
			"code":  "missing_url",
		})
		return
	}

	preview, err := h.metadataService.FetchPreview(c.Request.Context(), targetURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to fetch preview",
			"code":  "fetch_failed",
		})
		return
	}

	c.JSON(http.StatusOK, preview)
}
