package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

// Claims defines the structure of the JWT claims, including phone number and role
type Claims struct {
	UserID      int    `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token with the phone number and role
func GenerateJWT(userID int, phoneNumber, role string) (string, error) {
	claims := &Claims{
		UserID:      userID,
		PhoneNumber: phoneNumber,
		Role:        role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// VerifyJWT verifies the JWT token and extracts the phone number and role
func VerifyJWT(t string) (Claims, error) {
	claims := &Claims{}

	ts := strings.Split(t, "Bearer ")
	if len(ts) <= 1 {
		return Claims{}, errors.New("invalid or expired token")
	}
	tokenString := ts[1]

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Check that the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return Claims{}, errors.New("invalid or expired token")
	}

	// Return both the phone number and role
	return *claims, nil
}
