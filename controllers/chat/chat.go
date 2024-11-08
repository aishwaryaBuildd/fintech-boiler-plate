package chat

import (
	"fintech/store"
	"fintech/store/models"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	Store store.Store
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // For testing, allows any origin. Adjust this for production.
	},
}

// Connection represents a single WebSocket connection between two users.
type Connection struct {
	Conn       *websocket.Conn
	SenderID   int
	ReceiverID int
}

// A map to store active WebSocket connections.
// The key is a combination of senderID:receiverID.
var connections = make(map[string]*Connection)
var mutex sync.Mutex

func (controller ChatController) Chat(c *gin.Context) {
	senderID := c.MustGet("user_id").(int)
	courseID := c.Query("course_id")

	course, err := controller.Store.GetCourse(c, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	receiverID := course.AuthorID

	// Upgrade the connection to WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	conn := &Connection{
		Conn:       ws,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}
	key := fmt.Sprintf("%d:%d", senderID, receiverID)

	mutex.Lock()
	connections[key] = conn
	mutex.Unlock()

	fmt.Printf("New WebSocket connection: User %d to User %d\n", senderID, receiverID)

	// Listen for incoming messages
	go controller.handleMessages(c, conn)
}

func (controller ChatController) GetChatSessions(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	sessions, err := controller.Store.GetChatSessions(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (controller ChatController) GetChatSessionsMessages(c *gin.Context) {
	sessionID := c.Param("session_id")
	sID, _ := strconv.Atoi(sessionID)

	messages, err := controller.Store.GetChatSessionsMessages(c, sID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (controller ChatController) MarkChatSessionsAsRead(c *gin.Context) {
	sessionID := c.Param("session_id")
	sID, _ := strconv.Atoi(sessionID)

	err := controller.Store.MarkChatSessionsAsRead(c, sID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller ChatController) handleMessages(c *gin.Context, conn *Connection) {
	defer func() {
		conn.Conn.Close()
		key := fmt.Sprintf("%d:%d", conn.SenderID, conn.ReceiverID)
		mutex.Lock()
		delete(connections, key)
		mutex.Unlock()
	}()

	for {
		_, message, err := conn.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// Check and create sessions if needed
		m := models.Message{
			SenderID:   conn.SenderID,
			ReceiverID: conn.ReceiverID,
			Content:    string(message),
			CreatedAt:  time.Now(),
		}
		sessionID, err := controller.Store.GetOrCreateSession(c, m)
		if err != nil {
			fmt.Printf("User %d is not connected\n", conn.ReceiverID)
		}
		m.SessionID = sessionID
		err = controller.Store.AddMessage(c, m)
		if err != nil {
			fmt.Print("message not stored \n", m)
		}
		fmt.Printf("Message from User %d to User %d: %s\n", conn.SenderID, conn.ReceiverID, message)

		// Relay the message to the receiver if they're connected
		relayKey := fmt.Sprintf("%d:%d", conn.ReceiverID, conn.SenderID)
		mutex.Lock()
		receiverConn, exists := connections[relayKey]
		mutex.Unlock()

		if exists {
			receiverConn.Conn.WriteMessage(websocket.TextMessage, message)
		} else {
			fmt.Printf("User %d is not connected\n", conn.ReceiverID)
		}
	}
}
