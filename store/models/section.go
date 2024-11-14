package models

import (
	"time"

	"github.com/google/uuid"
)

type Section struct {
	ID          uuid.UUID `db:"id"`          // CHAR(36) UUID for folder ID
	Name        string    `db:"name"`        // VARCHAR(50), non-nullable
	Description string    `db:"description"` // VARCHAR(300), nullable
	CourseID    uuid.UUID `db:"course_id"`   // CHAR(36) UUID for course ID
	CreatedAt   time.Time `db:"created_at"`  // DATETIME(6) with default current timestamp
	UpdatedAt   time.Time `db:"updated_at"`  // DATETIME(6) with auto-update on current timestamp
}
