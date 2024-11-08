package auth

import (
	"database/sql"
	"errors"
	"fintech/store"
	"fintech/store/models"
	"fintech/utils"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Store store.Store
}

func (controller *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	adminPhoneNumbers := []string{"9840091130"}

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
	_, err := controller.Store.GetUserByPhoneNumber(c, user.PhoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error is %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}
	if err == nil {
		exists = 1
	}

	// Insert or update the OTP in the database
	if exists == 0 {
		// Example insertion query in Register function
		err := controller.Store.CreateUser(c, user.PhoneNumber, otp, otpExpiry, user.Role)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}
	} else {
		err := controller.Store.UpdateOTP(c, user.PhoneNumber, otp, otpExpiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP"})
			return
		}
	}

	// Send OTP via WhatsApp
	sendWhatsAppMessage(user.PhoneNumber, otp)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (controller *Controller) Verify(c *gin.Context) {
	var req VerifyRequest

	// Bind incoming JSON to the user model
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	u, err := controller.Store.GetUserByPhoneNumber(c, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Get stored OTP and its expiry time from the database

	// Check if the OTP matches and if it is not expired
	if u.OTP != req.OTP || time.Now().After(u.OTPExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Generate JWT token for the user with the role (default to "user" if role is not found)
	token, err := utils.GenerateJWT(u.ID, u.PhoneNumber, u.Role)
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

// Generate random OTP
func generateOTP() string {
	otp := rand.Intn(10000)         // Generate a random OTP
	return fmt.Sprintf("%04d", otp) // Format as a 4-digit string
}

// Send OTP via WhatsApp (dummy implementation, replace with actual sending logic)
func sendWhatsAppMessage(phone string, otp string) {
	fmt.Printf("Sending OTP %s to phone number %s via WhatsApp\n", otp, phone)
}
