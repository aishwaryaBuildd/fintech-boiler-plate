package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID `db:"id"`          // Matches CHAR(36) for UUID
	Name        string    `db:"name"`        // VARCHAR(50), non-nullable
	Description string    `db:"description"` // VARCHAR(300), nullable, use sql.NullString
	AuthorID    int       `db:"author_id"`   // INT, non-nullable
	FolderID    string    `db:"folder_id"`
	CreatedAt   time.Time `db:"created_at"` // DATETIME(6), default CURRENT_TIMESTAMP(6)
	UpdatedAt   time.Time `db:"updated_at"` // DATETIME(6), auto-updated with CURRENT_TIMESTAMP(6)
}
