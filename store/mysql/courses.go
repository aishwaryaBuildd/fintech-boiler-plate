package mysql

import (
	"context"
	"fintech/store/models"
)

func (m *MySQLStore) GetCourse(context context.Context, courseID string) (models.Course, error) {
	var c models.Course
	err := m.DB.GetContext(context, &c, "SELECT * FROM courses WHERE id = ?", courseID)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) ListCourse(context context.Context) ([]models.Course, error) {
	var c []models.Course
	err := m.DB.SelectContext(context, &c, "SELECT * FROM courses")
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) CreateCourse(context context.Context, c models.Course) error {
	_, err := m.DB.NamedExecContext(context, "INSERT INTO courses (id, name, description, author_id, folder_id, created_at, updated_at) VALUES (:id, :name, :description, :author_id, :folder_id, :created_at, :updated_at)",
		c)

	return err
}

func (m *MySQLStore) UpdateCourse(context context.Context, c models.Course) error {
	_, err := m.DB.NamedExecContext(context, "UPDATE courses SET name = :name, description = :description, author_id = :author_id, updated_at = :updated_at WHERE id = :id",
		c)
	return err
}

func (m *MySQLStore) DeleteCourse(context context.Context, id string) error {
	_, err := m.DB.ExecContext(context, "DELETE from courses WHERE id = ?",
		id)
	return err
}
