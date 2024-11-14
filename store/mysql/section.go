package mysql

import (
	"context"
	"fintech/store/models"
)

func (m *MySQLStore) GetSection(context context.Context, sectionID string) (models.Section, error) {
	var c models.Section
	err := m.DB.GetContext(context, &c, "SELECT * FROM sections WHERE id = ?", sectionID)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) ListSection(context context.Context) ([]models.Section, error) {
	var c []models.Section
	err := m.DB.SelectContext(context, &c, "SELECT * FROM sections")
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) CreateSection(context context.Context, c models.Section) error {
	_, err := m.DB.NamedExecContext(context, "INSERT INTO sections (id, name, description, course_id, section_id, created_at, updated_at) VALUES (:id, :name, :description, :course_id, :section_id, :created_at, :updated_at)",
		c)

	return err
}

func (m *MySQLStore) UpdateSection(context context.Context, c models.Section) error {
	_, err := m.DB.NamedExecContext(context, "UPDATE sections SET name = :name, description = :description, course_id = :course_id, updated_at = :updated_at WHERE id = :id",
		c)
	return err
}

func (m *MySQLStore) DeleteSection(context context.Context, id string) error {
	_, err := m.DB.ExecContext(context, "DELETE from sections WHERE id = ?",
		id)
	return err
}
