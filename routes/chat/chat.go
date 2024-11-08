package chat

import (
	chatController "fintech/controllers/chat"
	"fintech/utils"
	"net/http"

	"fintech/store"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine, db store.Store) {
	controller := chatController.ChatController{Store: db}

	r.GET("/chat/ws", ChatAuthMiddleware, controller.Chat)
	r.GET("/chat/sessions", ChatAuthMiddleware, controller.GetChatSessions)
	r.GET("/chat/sessions/:session_id/messages", ChatAuthMiddleware, controller.GetChatSessionsMessages)
	r.GET("/chat/sessions/:session_id/messages/read", ChatAuthMiddleware, controller.MarkChatSessionsAsRead)
}

func ChatAuthMiddleware(c *gin.Context) {
	token := c.Query("Authorization")
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
