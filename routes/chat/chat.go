package chat

import (
	chatController "fintech/controllers/chat"

	"fintech/middlewares"
	"fintech/store"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine, db store.Store) {
	controller := chatController.ChatController{Store: db}

	r.GET("/ws", middlewares.AuthMiddleware, controller.Chat)
}
