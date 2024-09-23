package mysql

import (
	"testing"

	"github.com/aishwaryaBuildd/go_boiler_plate/store"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStore with the mock database
	us := NewStore(db)

	user := store.User{
		Email:    "test@example.com",
		Password: "password",
	}
	// Expectations: Expect an INSERT query with specific arguments
	mock.ExpectExec("INSERT INTO users").
		WithArgs("test@example.com", "password").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test CreateUser method

	us.CreateUser(&user)
	assert.NoError(t, err)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStore with the mock database
	us := NewStore(db)

	// Expectations: Expect a SELECT query with a specific argument
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, "test@example.com", "password")
	mock.ExpectQuery("SELECT id, email, password FROM users").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	// Test GetUserByEmail method
	password, err := us.GetUserByEmail("test@example.com")
	assert.NoError(t, err)

	user := store.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	assert.Equal(t, &user, password)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
