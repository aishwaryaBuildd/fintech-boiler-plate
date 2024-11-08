package mysql

import (
	"context"
	"fintech/store/models"
)

func (m *MySQLStore) GetFolder(context context.Context, folderID string) (models.Folder, error) {
	var c models.Folder
	err := m.DB.GetContext(context, &c, "SELECT * FROM folders WHERE id = ?", folderID)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) ListFolder(context context.Context) ([]models.Folder, error) {
	var c []models.Folder
	err := m.DB.SelectContext(context, &c, "SELECT * FROM folders")
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *MySQLStore) CreateFolder(context context.Context, c models.Folder) error {
	_, err := m.DB.NamedExecContext(context, "INSERT INTO folders (id, name, description, course_id, folder_id, created_at, updated_at) VALUES (:id, :name, :description, :course_id, :folder_id, :created_at, :updated_at)",
		c)

	return err
}

func (m *MySQLStore) UpdateFolder(context context.Context, c models.Folder) error {
	_, err := m.DB.NamedExecContext(context, "UPDATE folders SET name = :name, description = :description, course_id = :course_id, updated_at = :updated_at WHERE id = :id",
		c)
	return err
}

func (m *MySQLStore) DeleteFolder(context context.Context, id string) error {
	_, err := m.DB.ExecContext(context, "DELETE from folders WHERE id = ?",
		id)
	return err
}
