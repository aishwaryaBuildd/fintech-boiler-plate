package models

import "time"

type ChatSession struct {
	ID          int       `db:"id"`           // ID of the session
	SenderID    int       `db:"sender_id"`    // User ID of the sender
	ReceiverID  int       `db:"receiver_id"`  // User ID of the receiver
	LastMessage string    `db:"last_message"` // Last message content in the session
	UnreadCount int       `db:"unread_count"` // Count of unread messages in the session
	Status      string    `db:"status"`       // Session status ('open' or 'closed')
	UpdatedAt   time.Time `db:"updated_at"`   // Timestamp of the last update
}
