package mysql

import (
	"database/sql"

	"github.com/aishwaryaBuildd/go_boiler_plate/store"
)

type MySQLUserStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *MySQLUserStore {
	return &MySQLUserStore{db: db}
}

func (s *MySQLUserStore) CreateUser(user *store.User) error {
	_, err := s.db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
	return err
}

func (s *MySQLUserStore) GetUserByEmail(email string) (*store.User, error) {
	row := s.db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email)
	user := &store.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewDB() (*sql.DB, error) {
	return sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
}
