package controllers

import (
	"net/http"
	"os"

	"messaging-platform/database"
	"messaging-platform/models"

	"github.com/gin-gonic/gin"
)

// HealthCheck - Check if the server and database are running
func HealthCheck(c *gin.Context) {
	sqlDB, err := database.DB.DB() // Get the raw database connection
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to connect to the database"})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database is unreachable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Server and database are running"})
}

// FileUpload - Upload attachments for messages
func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// Define the upload directory
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}

	// Save the file
	filePath := uploadDir + "/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": filePath})
}

// SearchUsersOrGroups - Search for users or groups by name or other filters
func SearchUsersOrGroups(c *gin.Context) {
	query := c.Query("q")         // Get the search query
	searchType := c.Query("type") // Either "users" or "groups"

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	if searchType == "users" {
		var users []models.User
		if err := database.DB.Where("username LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
			return
		}
		c.JSON(http.StatusOK, users)
	} else if searchType == "groups" {
		var groups []models.Group
		if err := database.DB.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&groups).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search groups"})
			return
		}
		c.JSON(http.StatusOK, groups)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search type"})
	}
}
