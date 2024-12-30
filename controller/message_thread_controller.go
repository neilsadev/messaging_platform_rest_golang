package controllers

import (
	"net/http"
	"time"

	"messaging-platform/database"
	"messaging-platform/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateThread - Reply to a specific message, creating a thread
func CreateThread(c *gin.Context) {
	var thread models.MessageThread
	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate parent message
	var parentMessage models.Message
	if err := database.DB.First(&parentMessage, "message_uuid = ?", thread.ParentMessageUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent message not found"})
		return
	}

	// Validate child message
	var childMessage models.Message
	if err := database.DB.First(&childMessage, "message_uuid = ?", thread.ChildMessageUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child message not found"})
		return
	}

	// Generate a new UUID for the thread
	thread.ThreadUUID = uuid.New().String()
	thread.CreatedAt = time.Now()

	// Save the thread
	if err := database.DB.Create(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thread"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Thread created successfully", "data": thread})
}

// GetThread - Fetch all messages in a specific thread
func GetThread(c *gin.Context) {
	parentMessageUUID := c.Param("parentMessageUUID")

	var threadMessages []models.Message
	err := database.DB.Joins("JOIN message_threads ON message_threads.child_message_uuid = messages.message_uuid").
		Where("message_threads.parent_message_uuid = ?", parentMessageUUID).
		Find(&threadMessages).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch thread messages"})
		return
	}

	c.JSON(http.StatusOK, threadMessages)
}
