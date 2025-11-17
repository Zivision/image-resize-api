package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Set lower memory for multipart forms
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Health check endpoint (unversioned, for AWS load balancer)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 endpoints
	{
		v1 := r.Group("/api/v1")
		v1.GET("/test", testEndpoint)
		v1.POST("/image", imageEndpoint)
	}

	// Get port from environment variable (AWS sets this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	// Listen for SIGINT (Ctrl+C) and SIGTERM (AWS/Docker stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	// Give in-flight requests up to 10 seconds to complete
	// (10 seconds is good for image processing - adjust if needed)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited cleanly")
}
