package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Set lower memory for multipart forms
	r.MaxMultipartMemory = 8 << 20 // 8 Mib

	// API v1 endpoints
	{
		v1 := r.Group("/api/v1")
		v1.GET("/test", testEndpoint)
		v1.POST("/image", imageEndpoint)
	}
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
