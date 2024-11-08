package mysql

import (
	"context"
	"fintech/store/models"
)

func (m *MySQLStore) AddMessage(context context.Context, c models.Message) error {
	_, err := m.DB.NamedExecContext(context, "INSERT INTO messages (id, sender_id, receiver_id, content, created_at) VALUES (:id, :sender_id, :receiver_id, :content, :created_at)",
		c)

	return err
}

func (m *MySQLStore) GetOrCreateSession(context context.Context, c models.Message) (int, error) {
	var sessionID int
	err := m.DB.QueryRowContext(context, `
        INSERT INTO chat_sessions (sender_id, receiver_id, last_message, unread_count)
        VALUES ($1, $2, $3, 1)
        ON CONFLICT (sender_id, receiver_id) 
        DO UPDATE SET last_message = $3, unread_count = chat_sessions.unread_count + 1, updated_at = NOW()
        RETURNING id`,
		c.SenderID, c.ReceiverID, c.Content).Scan(&sessionID)

	if err != nil {
		return sessionID, err
	}

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
	_, err := m.DB.NamedExecContext(context, "UPDATE chat_sessions SET unread_count = 0 WHERE id = ?",
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
