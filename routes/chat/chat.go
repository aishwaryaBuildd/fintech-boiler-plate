package chat

import (
	chatController "fintech/controllers/chat"
	"fintech/middlewares"
	"fintech/utils"
	"net/http"

	"fintech/store"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine, db store.Store) {
	controller := chatController.ChatController{Store: db}

	r.GET("/chat/ws", ChatAuthMiddleware, controller.Chat)
	r.GET("/chat/ws/admin", ChatAdminMiddleware, controller.ChatAdmin)
	r.GET("/chat/sessions", middlewares.AuthMiddleware, controller.GetChatSessions)
	r.GET("/chat/sessions/:session_id/messages", middlewares.AuthMiddleware, controller.GetChatSessionsMessages)
	r.GET("/chat/sessions/:session_id/messages/read", middlewares.AuthMiddleware, controller.MarkChatSessionsAsRead)
}

func ChatAuthMiddleware(c *gin.Context) {
	token := c.Query("authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	claims, err := utils.VerifyJWT(token) // We only care about phone number here
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("user_id", claims.UserID)           // Store phone number in context
	c.Set("phone_number", claims.PhoneNumber) // Store phone number in context
	c.Next()
}

// AdminMiddleware checks if the user is an admin
func ChatAdminMiddleware(c *gin.Context) {
	token := c.Query("authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	claims, err := utils.VerifyJWT(token) // We only care about phone number here
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access"})
		c.Abort()
		return
	}
	c.Set("user_id", claims.UserID)           // Store phone number in context
	c.Set("phone_number", claims.PhoneNumber) // Store phone number in context
	c.Next()
}
