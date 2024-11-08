package models

import "time"

type Message struct {
	ID         int       `db:"id"`          // SERIAL primary key
	SenderID   int       `db:"sender_id"`   // INT, foreign key referencing users(id)
	ReceiverID int       `db:"receiver_id"` // INT, foreign key referencing users(id)
	Content    string    `db:"content"`     // TEXT, non-nullable message content
	CreatedAt  time.Time `db:"created_at"`  // TIMESTAMP with default current timestamp
}
