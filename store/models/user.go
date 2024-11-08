package models

import "time"

// User represents the user model in the application
type User struct {
	ID          int       `db:"id"`           // Unique identifier for the user
	PhoneNumber string    `db:"phone_number"` // User's phone number
	OTP         string    `db:"otp_code"`     // One-time password for verification
	OTPExpiry   time.Time `db:"otp_expiry"`   // Expiration time for the OTP
	Role        string    `db:"role"`         // Role of the user (e.g., "user", "admin")
	CreatedAt   time.Time `db:"created_at"`   // Timestamp of user creation
	UpdatedAt   time.Time `db:"updated_at"`   // Timestamp of the last update
}
