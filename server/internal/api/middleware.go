package api

import (
	"fmt"
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

		// Log to console, not to response
		// This is the key fix - use fmt.Fprintf with gin.DefaultWriter instead of c.String
		fmt.Fprintf(gin.DefaultWriter, "REQUEST|%s|%s|%d|%v ms\n",
			reqMethod,
			reqURI,
			statusCode,
			latencyTime.Milliseconds(),
		)
	}
}
