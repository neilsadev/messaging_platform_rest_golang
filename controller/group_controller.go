package controllers

import (
	"net/http"

	"messaging-platform/database"
	"messaging-platform/models"

	"github.com/gin-gonic/gin"
)

// GetGroups retrieves all groups
func GetGroups(c *gin.Context) {
	var groups []models.Group
	if err := database.DB.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// GetGroup retrieves a single group by ID
func GetGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group
	if err := database.DB.First(&group, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
}

// CreateGroup creates a new group
func CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}
	c.JSON(http.StatusCreated, group)
}

// UpdateGroup updates an existing group by ID
func UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group

	// Check if the group exists
	if err := database.DB.First(&group, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Bind the new data
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save updates
	if err := database.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group"})
		return
	}
	c.JSON(http.StatusOK, group)
}

// DeleteGroup deletes a group by ID
func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Group{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}
