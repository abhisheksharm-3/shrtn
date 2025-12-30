// Package middleware provides HTTP middleware for the API server.
package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

const apiKeyHeader = "X-API-Key"

// Security returns middleware that adds security headers to responses.
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// APIKeyAuth returns middleware that validates API key for protected routes.
func APIKeyAuth(validAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validAPIKey == "" {
			c.Next()
			return
		}

		providedKey := c.GetHeader(apiKeyHeader)
		if providedKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
				"code":  "missing_api_key",
			})
			return
		}

		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(validAPIKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
				"code":  "invalid_api_key",
			})
			return
		}

		c.Next()
	}
}
