// controller/auth_controller.go

package controller

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"fintech/models"
	"fintech/services"
	"fintech/utils"

	"github.com/gin-gonic/gin"
)

// Generate random OTP
func generateOTP() string {
	otp := rand.Intn(10000)         // Generate a random OTP
	return fmt.Sprintf("%04d", otp) // Format as a 4-digit string
}

// Send OTP via WhatsApp (dummy implementation, replace with actual sending logic)
func sendWhatsAppMessage(phone string, otp string) {
	fmt.Printf("Sending OTP %s to phone number %s via WhatsApp\n", otp, phone)
}

// Register handles user registration and OTP generation
func Register(c *gin.Context, db *sql.DB) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	adminPhoneNumbers := []string{}

	role := "user"

	if slices.Contains(adminPhoneNumbers, req.PhoneNumber) {
		role = "admin"
	}

	user := models.User{
		PhoneNumber: req.PhoneNumber,
		Role:        role,
	}

	otp := generateOTP()                            // Generate a random OTP
	otpExpiry := time.Now().Add(1440 * time.Minute) // Set OTP expiry time

	// Check if user already exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE phone_number = ?", user.PhoneNumber).Scan(&exists)
	if err != nil {
		log.Printf("error is %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	// Insert or update the OTP in the database
	if exists == 0 {
		// Example insertion query in Register function
		_, err := db.Exec("INSERT INTO users (phone_number, otp_code, otp_expiry, role) VALUES (?, ?, ?, ?)",
			user.PhoneNumber, otp, otpExpiry, user.Role)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}
	} else {
		_, err := db.Exec("UPDATE users SET otp_code = ?, otp_expiry = ? WHERE phone_number = ?",
			otp, otpExpiry, user.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP"})
			return
		}
	}

	// Send OTP via WhatsApp
	sendWhatsAppMessage(user.PhoneNumber, otp)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Verify handles user OTP verification
func Verify(c *gin.Context, db *sql.DB) {
	var req VerifyRequest

	// Bind incoming JSON to the user model
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get stored OTP and its expiry time from the database
	storedOTP, otpExpiry, err := services.GetUserOTP(req.PhoneNumber, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Check if the OTP matches and if it is not expired
	if storedOTP != req.OTP || time.Now().After(otpExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Fetch the user's role from the database in the same function
	var role string
	err = db.QueryRow("SELECT role FROM users WHERE phone_number = ?", req.PhoneNumber).Scan(&role)
	if err != nil || role == "" {
		// If role is not found or an error occurs, default to "user"
		role = "user"
	}

	// Generate JWT token for the user with the role (default to "user" if role is not found)
	token, err := utils.GenerateJWT(req.PhoneNumber, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the JWT token in the response
	c.JSON(http.StatusOK, gin.H{"token": token})
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" validate:"min=10,max=10"`
}

type VerifyRequest struct {
	PhoneNumber string `json:"phone_number" validate:"min=10,max=10"`
	OTP         string `json:"otp" validate:"min=4,max=4"`
}
