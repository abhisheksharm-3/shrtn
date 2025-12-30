// Package middleware provides HTTP middleware for the API server.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns middleware that logs request details.
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return param.TimeStamp.Format(time.RFC3339) +
				" | " + param.Method +
				" | " + param.Path +
				" | " + param.StatusCodeColor() + " " + 
				string(rune(param.StatusCode)) + param.ResetColor() +
				" | " + param.Latency.String() +
				" | " + param.ClientIP +
				"\n"
		},
		Output:    gin.DefaultWriter,
		SkipPaths: []string{"/health"},
	})
}
