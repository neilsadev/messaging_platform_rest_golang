package controllers

import (
	"net/http"
	"time"

	"messaging-platform/database"
	"messaging-platform/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SendMessage - Send a new message to a user or group
func SendMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new UUID for the message
	message.MessageUUID = uuid.New().String()
	message.CreatedAt = time.Now()

	// Validate sender
	var sender models.User
	if err := database.DB.First(&sender, "user_uuid = ?", message.SenderUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender not found"})
		return
	}

	// If it's a group message, validate the group
	if message.GroupUUID != nil {
		var group models.Group
		if err := database.DB.First(&group, "group_uuid = ?", *message.GroupUUID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
	}

	// Save the message
	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent successfully", "data": message})
}

// GetMessage - Fetch details of a specific message by ID
func GetMessage(c *gin.Context) {
	messageUUID := c.Param("messageUUID")

	var message models.Message
	if err := database.DB.First(&message, "message_uuid = ?", messageUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}

// UpdateMessage - Edit the content of an existing message
func UpdateMessage(c *gin.Context) {
	messageUUID := c.Param("messageUUID")

	var message models.Message
	if err := database.DB.First(&message, "message_uuid = ?", messageUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.Content = input.Content
	message.UpdatedAt = time.Now()

	if err := database.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully", "data": message})
}

// DeleteMessage - Soft delete a message
func DeleteMessage(c *gin.Context) {
	messageUUID := c.Param("messageUUID")

	var message models.Message
	if err := database.DB.First(&message, "message_uuid = ?", messageUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	now := time.Now()
	message.DeletedAt = &now

	if err := database.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// ListMessages - Retrieve messages in a direct conversation or group
func ListMessages(c *gin.Context) {
	var messages []models.Message

	groupUUID := c.Query("group_uuid")
	senderUUID := c.Query("sender_uuid")
	recipientUUID := c.Query("recipient_uuid")

	query := database.DB
	if groupUUID != "" {
		query = query.Where("group_uuid = ?", groupUUID)
	} else if senderUUID != "" && recipientUUID != "" {
		query = query.Where(
			"(sender_uuid = ? AND recipient_uuid = ?) OR (sender_uuid = ? AND recipient_uuid = ?)",
			senderUUID, recipientUUID, recipientUUID, senderUUID,
		)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	if err := query.Order("created_at ASC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
