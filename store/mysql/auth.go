package mysql

import (
	"context"
	"fintech/store/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MySQLStore struct {
	DB *sqlx.DB
}

func NewMySQLStore(db *sqlx.DB) *MySQLStore {
	return &MySQLStore{
		DB: db,
	}
}

func (m *MySQLStore) GetUserByPhoneNumber(context context.Context, phoneNumber string) (models.User, error) {
	var u models.User
	err := m.DB.GetContext(context, &u, "SELECT * FROM users WHERE phone_number = ?", phoneNumber)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *MySQLStore) CreateUser(context context.Context, phoneNumber, otp string, otpExpiry time.Time, role string) error {
	_, err := m.DB.ExecContext(context, "INSERT INTO users (phone_number, otp_code, otp_expiry, role) VALUES (?, ?, ?, ?)",
		phoneNumber, otp, otpExpiry, role)

	return err
}

func (m *MySQLStore) UpdateOTP(context context.Context, phoneNumber, otp string, otpExpiry time.Time) error {
	_, err := m.DB.ExecContext(context, "UPDATE users SET otp_code = ?, otp_expiry = ? WHERE phone_number = ?",
		otp, otpExpiry, phoneNumber)
	return err
}
