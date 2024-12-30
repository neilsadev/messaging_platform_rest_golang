package controllers

import (
	"net/http"
	"time"

	"messaging-platform/database"
	"messaging-platform/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddUserToGroup adds a user to a group
func AddUserToGroup(c *gin.Context) {
	groupUUID := c.Param("groupUUID")

	var membership models.GroupMembership
	if err := c.ShouldBindJSON(&membership); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	membership.GroupUUID = groupUUID
	membership.MembershipUUID = GenerateUUID()
	membership.JoinedAt = time.Now()

	// Validate Group and User
	var group models.Group
	if err := database.DB.First(&group, "group_uuid = ?", groupUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "user_uuid = ?", membership.UserUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Save the membership
	if err := database.DB.Create(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to group"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added to group", "membership": membership})
}

// ListGroupMembers lists all users in a specific group
func ListGroupMembers(c *gin.Context) {
	groupUUID := c.Param("groupUUID")

	var members []models.User
	err := database.DB.
		Joins("JOIN group_memberships ON group_memberships.user_uuid = users.user_uuid").
		Where("group_memberships.group_uuid = ?", groupUUID).
		Find(&members).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

// RemoveUserFromGroup removes a user from a group
func RemoveUserFromGroup(c *gin.Context) {
	groupUUID := c.Param("groupUUID")
	userUUID := c.Param("userUUID")

	if err := database.DB.Delete(&models.GroupMembership{}, "group_uuid = ? AND user_uuid = ?", groupUUID, userUUID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from group"})
}

// GenerateUUID generates a unique identifier for the membership
func GenerateUUID() string {
	return uuid.New().String()
}
