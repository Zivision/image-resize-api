package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/davidbyttow/govips/v2/vips"
)

/*
 *
 *  Util Functions
 *
 */
func processImage(imageData []byte, width int, height int) ([]byte, error) {
	// Loads imageData ([]bytes) and returns imageRef (*vips.ImageRef)
	imageRef, err := vips.NewImageFromBuffer(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to load image from buffer: %w", err)
	}
	defer imageRef.Close() // Close image ref after completion

	err = imageRef.Thumbnail(width, height, vips.Interesting(vips.InterestingAll))
	if err != nil {
		return nil, fmt.Errorf("failed to thumbnail image: %w", err)
	}

	exportParams := vips.NewJpegExportParams()
	exportParams.Quality = 85 // Set desired quality
	exportParams.StripMetadata = true // Remove unnecessary metadata

	imageOutputBytes, _, err := imageRef.ExportJpeg(exportParams)
	if err != nil {
		return nil, fmt.Errorf("failed to export image to bytes: %w", err)
	}

	return imageOutputBytes, nil
}

func processFileContents(data []byte) ([]byte, error) {
	// --- VIPS Initialization ---
	// This is needed to call the function above
	// MUST be called before any other VIPS function.
	vips.Startup(&vips.Config{
		ConcurrencyLevel: 1, // Control internal threading
	})
	defer vips.Shutdown() // MUST be called to clean up resources


	header := []byte("--- Processed by API ---\n")

	// Create a new byte slice by combining the header and the original data
	var processedData bytes.Buffer
	processedData.Write(header)
	processedData.Write(data)

	resizedImage, err := processImage(processedData.Bytes(), 500, 500)
	if err != nil {
		return nil, fmt.Errorf("Failed to process image bytes: %w", err)
	}

	return resizedImage, nil

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

func fileEndpoint(c *gin.Context) {
	// Takes file from endpoint
	fileHeader, err := c.FormFile("file") // "file" is expected form field key
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve file: " + err.Error(),
		})
		return
	}

	// Attempts to open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file: " + err.Error()})
		return
	}
	// Close file
	defer file.Close()

	// Process file into bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file contents: " + err.Error()})
		return
	}

	// Process file contents
	processedData, err := processFileContents(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image process failed: " + err.Error()})
		return
	}

	// Set content type
	contentType := "image/jpeg"

	c.Data(http.StatusOK, contentType, processedData)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"processed_%s\"", fileHeader.Filename))
}
