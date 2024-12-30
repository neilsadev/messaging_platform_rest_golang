package controllers

import (
	"net/http"
	"time"

	"messaging-platform/database"
	"messaging-platform/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddReaction - React to a message
func AddReaction(c *gin.Context) {
	var reaction models.Reaction
	if err := c.ShouldBindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new UUID for the reaction
	reaction.ReactionUUID = uuid.New().String()
	reaction.CreatedAt = time.Now()

	// Validate message
	var message models.Message
	if err := database.DB.First(&message, "message_uuid = ?", reaction.MessageUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Validate user
	var user models.User
	if err := database.DB.First(&user, "user_uuid = ?", reaction.UserUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Save the reaction
	if err := database.DB.Create(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reaction added successfully", "data": reaction})
}

// RemoveReaction - Remove a reaction from a message
func RemoveReaction(c *gin.Context) {
	reactionUUID := c.Param("reactionUUID")

	var reaction models.Reaction
	if err := database.DB.First(&reaction, "reaction_uuid = ?", reactionUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reaction not found"})
		return
	}

	if err := database.DB.Delete(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction removed successfully"})
}

// ListReactions - Retrieve all reactions for a specific message
func ListReactions(c *gin.Context) {
	messageUUID := c.Param("messageUUID")

	var reactions []models.Reaction
	if err := database.DB.Where("message_uuid = ?", messageUUID).Find(&reactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reactions"})
		return
	}

	c.JSON(http.StatusOK, reactions)
}
