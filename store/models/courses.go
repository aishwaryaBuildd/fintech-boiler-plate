package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID `db:"id"`          // CHAR(36) for UUID, primary key
	Name        string    `db:"name"`        // VARCHAR(50), non-nullable
	Description string    `db:"description"` // VARCHAR(300), nullable, use sql.NullString
	Level       string    `db:"level"`       // VARCHAR(50), non-nullable
	Category    string    `db:"category"`    // VARCHAR(100), non-nullable
	Topic       string    `db:"topic"`       // VARCHAR(100), non-nullable
	Duration    string    `db:"duration"`    // VARCHAR(50), non-nullable
	AuthorID    int       `db:"author_id"`   // INT, non-nullable
	FolderID    string    `db:"folder_id"`   // VARCHAR(200), non-nullable
	Thumbnail   string    `db:"thumbnail"`   // VARCHAR(300), non-nullable
	Trailer     string    `db:"trailer"`     // VARCHAR(300), nullable, use sql.NullString
	CreatedAt   time.Time `db:"created_at"`  // DATETIME(6), default CURRENT_TIMESTAMP(6)
	UpdatedAt   time.Time `db:"updated_at"`  // DATETIME(6), auto-updated with CURRENT_TIMESTAMP(6)
}
