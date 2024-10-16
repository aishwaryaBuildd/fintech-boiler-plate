package models

// User represents a user in the system

import (
	"time"
)

// User represents the user model in the application
type User struct {
	ID          int       `json:"id"`           // Unique identifier for the user
	PhoneNumber string    `json:"phone_number"` // User's phone number
	OTP         string    `json:"otp"`          // One-time password for verification
	OTPExpiry   time.Time `json:"otp_expiry"`   // Expiration time for the OTP
	Role        string    `json:"role"`         // Role of the user (e.g., "user", "admin")
	CreatedAt   time.Time `json:"created_at"`   // Timestamp of user creation
	UpdatedAt   time.Time `json:"updated_at"`   // Timestamp of the last update
}
