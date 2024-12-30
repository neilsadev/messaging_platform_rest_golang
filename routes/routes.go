package routes

import (
	controllers "messaging-platform/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(routes *gin.Engine) {
	// User Routes
	userRoutes := routes.Group("/users")
	{
		userRoutes.POST("/", controllers.CreateUser)      // Create User
		userRoutes.GET("/", controllers.ListUsers)        // List Users
		userRoutes.GET("/:id", controllers.GetUser)       // Get User by ID
		userRoutes.PUT("/:id", controllers.UpdateUser)    // Update User
		userRoutes.DELETE("/:id", controllers.DeleteUser) // Delete User
	}

	// Auth Routes
	authRoutes := routes.Group("/auth")
	{
		authRoutes.POST("/login", controllers.UserLogin)      // User login
		authRoutes.POST("/logout", controllers.UserLogout)    // User logout
		authRoutes.POST("/refresh", controllers.RefreshToken) // Refresh token
	}

	// Group Routes
	groupRoutes := routes.Group("/groups")
	{
		groupRoutes.GET("/", controllers.GetGroups)         // List all groups
		groupRoutes.GET("/:id", controllers.GetGroup)       // Get a specific group by ID
		groupRoutes.POST("/", controllers.CreateGroup)      // Create a new group
		groupRoutes.PUT("/:id", controllers.UpdateGroup)    // Update an existing group by ID
		groupRoutes.DELETE("/:id", controllers.DeleteGroup) // Delete a group by ID

		// Member-specific routes under /groups/:id/members
		groupRoutes.POST("/:id/members", controllers.AddUserToGroup)                  // Add a user to a group
		groupRoutes.GET("/:id/members", controllers.ListGroupMembers)                 // List all members of a group
		groupRoutes.DELETE("/:id/members/:userUUID", controllers.RemoveUserFromGroup) // Remove a user from a group
	}

	// Message Routes
	messageRoutes := routes.Group("/messages")
	{
		messageRoutes.POST("/", controllers.SendMessage)                 // Send a new message
		messageRoutes.GET("/:messageUUID", controllers.GetMessage)       // Get a specific message by ID
		messageRoutes.PUT("/:messageUUID", controllers.UpdateMessage)    // Update a message
		messageRoutes.DELETE("/:messageUUID", controllers.DeleteMessage) // Delete a message
		messageRoutes.GET("/", controllers.ListMessages)                 // List messages (direct or group)
	}

	// Thread Routes
	threadRoutes := routes.Group("/threads")
	{
		threadRoutes.POST("/", controllers.CreateThread)               // Create a new thread (reply to a message)
		threadRoutes.GET("/:parentMessageUUID", controllers.GetThread) // Get all messages in a thread
	}

	// Reaction Routes
	reactionRoutes := routes.Group("/reactions")
	{
		reactionRoutes.POST("/", controllers.AddReaction)                   // Add a reaction to a message
		reactionRoutes.DELETE("/:reactionUUID", controllers.RemoveReaction) // Remove a reaction
		reactionRoutes.GET("/:messageUUID", controllers.ListReactions)      // List all reactions for a message
	}

	// Utility Routes
	utilityRoutes := routes.Group("/utility")
	{
		utilityRoutes.GET("/health", controllers.HealthCheck)         // Health check
		utilityRoutes.POST("/file-upload", controllers.FileUpload)    // File upload
		utilityRoutes.GET("/search", controllers.SearchUsersOrGroups) // Search users or groups
	}
}
