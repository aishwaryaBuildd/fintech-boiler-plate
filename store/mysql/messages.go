package mysql

import (
	"context"
	"fintech/store/models"
)

func (m *MySQLStore) AddMessage(context context.Context, c models.Message) error {
	_, err := m.DB.NamedExecContext(context, "INSERT INTO messages (id, sender_id, receiver_id, content, session_id,created_at) VALUES (:id, :sender_id, :receiver_id, :content, :session_id, :created_at)",
		c)

	return err
}

func (m *MySQLStore) GetOrCreateSession(context context.Context, c models.Message) (int, error) {
	var sessionID int
	query := `
        INSERT INTO chat_sessions (sender_id, receiver_id, last_message, unread_count, updated_at)
        VALUES (?, ?, ?, 1, NOW())
        ON DUPLICATE KEY UPDATE 
            last_message = VALUES(last_message), 
            unread_count = chat_sessions.unread_count + 1, 
            updated_at = NOW()`

	result, err := m.DB.ExecContext(context, query, c.SenderID, c.ReceiverID, c.Content)
	if err != nil {
		return 0, err
	}

	// Get the session ID
	sessionID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	sessionID = int(sessionID64)

	return sessionID, nil
}

func (m *MySQLStore) GetChatSessions(context context.Context, userID int) ([]models.ChatSession, error) {
	var c []models.ChatSession
	err := m.DB.SelectContext(context, &c, "SELECT * FROM chat_sessions where sender_id = ? OR receiver_id = ?", userID, userID)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) MarkChatSessionsAsRead(context context.Context, ChatSessionID int) error {
	_, err := m.DB.ExecContext(context, "UPDATE messages SET is_read = 1 WHERE session_id = ?",
		ChatSessionID)

	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(context, "UPDATE chat_sessions SET unread_count = 0 WHERE id = ?",
		ChatSessionID)
	return err
}

func (m *MySQLStore) GetChatSessionsMessages(context context.Context, sessionID int) ([]models.Message, error) {
	var c []models.Message
	err := m.DB.SelectContext(context, &c, "SELECT * FROM messages where session_id = ?", sessionID)
	if err != nil {
		return c, err
	}

	return c, nil
}
