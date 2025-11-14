package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

/*
 *
 *  Util Functions
 *
 */
func processJpeg(jpegBytes []byte) ([]byte, error) {
	// Deocode image bytes
	img, err := jpeg.Decode(bytes.NewReader(jpegBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image bytes: %w", err)
	}

	// Flip Vertically (Will add more features later)
	flippedImg := imaging.FlipV(img)

	// Declare return byte buffer
	var buf bytes.Buffer

	// Re-encode it
	err = jpeg.Encode(&buf, flippedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		return nil, fmt.Errorf("failed to encode flipped image: %w", err)
	}

	// Return result
	return buf.Bytes(), nil
}

func sortImageType(data []byte) ([]byte, error) {
	// TODO Will add sorting for different file types and operations
	proccesedImg, err := processJpeg(data)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}
	return proccesedImg, nil
}

/*
 *
 * Endpoint functions
 *
 */

func testEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Server is running",
		"status":  http.StatusOK,
	})
}

func imageEndpoint(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded",
		})
		return
	}

	if file.Size > 10*1024*1024 { // 10MB limit
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 10MB",
		})
		return
	}

	// Attempts to open the file
	image, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file: " + err.Error()})
		return
	}
	// Close file
	defer image.Close()

	// Process image into bytes
	imageBytes, err := io.ReadAll(image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file contents: " + err.Error()})
		return
	}

	// Process file contents
	processedData, err := sortImageType(imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image process failed: " + err.Error()})
		return
	}

	// Set content type
	contentType := "image/jpeg"

	// Send image back
	c.Data(http.StatusOK, contentType, processedData)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"processed_%s\"", file.Filename))
}
