package models

import (
	"time"

	"github.com/google/uuid"
)

type Lecture struct {
	ID           uuid.UUID `db:"id"`            // CHAR(36), non-nullable
	Name         string    `db:"name"`          // VARCHAR(50), non-nullable
	Description  string    `db:"description"`   // VARCHAR(300), nullable
	SectionID    uuid.UUID `db:"section_id"`    // CHAR(36), non-nullable
	VideoID      string    `db:"video_id"`      // VARCHAR(50), nullable (assuming some lectures may not have a video)
	LectureNotes string    `db:"lecture_notes"` // VARCHAR(1000), nullable
	FileName     string    `db:"file_name"`     // VARCHAR(50), nullable
	FileURL      string    `db:"file_url"`      // VARCHAR(50), nullable
	CreatedAt    time.Time `db:"created_at"`    // DATETIME(6), default CURRENT_TIMESTAMP(6)
	UpdatedAt    time.Time `db:"updated_at"`    // DATETIME(6), auto-updated with CURRENT_TIMESTAMP(6)
}
