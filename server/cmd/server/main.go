package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a default Gin router
	router := gin.Default()

	// Handler for the root route
	router.GET("/", func(c *gin.Context) {
		c.String(200, "URL Shortener Mock Server is running! Now in Docker too! yay! yhis is a test. i will now do a full test")
	})

	// Start the server on port 8080
	port := ":8080"
	log.Printf("Server starting on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
