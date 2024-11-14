package mysql

import (
	"context"
	"fintech/store/models"
)

func (m *MySQLStore) GetLecture(context context.Context, lectureID string) (models.Lecture, error) {
	var c models.Lecture
	err := m.DB.GetContext(context, &c, "SELECT * FROM lectures WHERE id = ?", lectureID)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) ListLecture(context context.Context) ([]models.Lecture, error) {
	var c []models.Lecture
	err := m.DB.SelectContext(context, &c, "SELECT * FROM lectures")
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) CreateLecture(context context.Context, c models.Lecture) error {
	_, err := m.DB.NamedExecContext(context, `
		INSERT INTO lectures (
			id, name, description, section_id, video_id, lecture_notes, file_name, file_url, created_at, updated_at
		) VALUES (
			:id, :name, :description, :section_id, :video_id, :lecture_notes, :file_name, :file_url, :created_at, :updated_at
		)`, c)

	return err
}

func (m *MySQLStore) UpdateLecture(context context.Context, c models.Lecture) error {
	_, err := m.DB.NamedExecContext(context, `
		UPDATE lectures 
		SET name = :name, 
		    description = :description, 
		    section_id = :section_id, 
		    video_id = :video_id, 
		    lecture_notes = :lecture_notes, 
		    file_name = :file_name, 
		    file_url = :file_url, 
		    updated_at = :updated_at 
		WHERE id = :id`, c)

	return err
}

func (m *MySQLStore) DeleteLecture(context context.Context, id string) error {
	_, err := m.DB.ExecContext(context, "DELETE from lectures WHERE id = ?",
		id)
	return err
}
