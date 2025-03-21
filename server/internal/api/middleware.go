package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs incoming requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Starting time
		startTime := time.Now()

		// Process request
		c.Next()

		// End time
		endTime := time.Now()

		// Execution time
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := c.Request.Method

		// Request route
		reqURI := c.Request.RequestURI

		// Status code
		statusCode := c.Writer.Status()

		// Log request details
		// Using Gin's logging capabilities
		c.Set("request_time", latencyTime)
		c.Next()
		
		// Log format: REQUEST | method | uri | status | latency
		// For example: REQUEST | GET | /api/shorten | 200 | 34.3ms
		c.Writer.Status()
		c.String(statusCode, "REQUEST | %s | %s | %d | %v\n", 
			reqMethod, 
			reqURI, 
			statusCode, 
			latencyTime,
		)
	}
}
