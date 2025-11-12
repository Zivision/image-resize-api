package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 *
 *  Util Functions
 *
 */

func processFileContents(data []byte) []byte {
	// For this example, we'll just prepend a simple header to simulate processing.
	header := []byte("--- Processed by API ---\n")

	// Create a new byte slice by combining the header and the original data
	var processedData bytes.Buffer
	processedData.Write(header)
	processedData.Write(data)

	return processedData.Bytes()
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
	processedData := processFileContents(fileBytes)

	// Set content type
	contentType := "image/png"

	c.Data(http.StatusOK, contentType, processedData)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"processed_%s\"", fileHeader.Filename))
}
