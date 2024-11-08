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
