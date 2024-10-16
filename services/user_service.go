package services

import (
	"database/sql"
	"fmt"
	"time"
)

func GetUserOTP(phoneNumber string, db *sql.DB) (string, time.Time, error) {
	var otpCode string
	var otpExpiryRaw []uint8 // Use a []uint8 to capture raw data first

	err := db.QueryRow("SELECT otp_code, otp_expiry FROM users WHERE phone_number = ?", phoneNumber).Scan(&otpCode, &otpExpiryRaw)
	if err != nil {
		return "", time.Time{}, err
	}

	// Convert the []uint8 (which is a byte slice) to a string
	otpExpiryStr := string(otpExpiryRaw)

	// Parse the string to time.Time
	otpExpiry, err := time.Parse("2006-01-02 15:04:05", otpExpiryStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse otp_expiry: %v", err)
	}

	return otpCode, otpExpiry, nil
}
