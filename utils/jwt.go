package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

// Claims defines the structure of the JWT claims, including phone number and role
type Claims struct {
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token with the phone number and role
func GenerateJWT(phoneNumber, role string) (string, error) {
	claims := &Claims{
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
func VerifyJWT(tokenString string) (string, string, error) {
	claims := &Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Check that the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid or expired token")
	}

	// Return both the phone number and role
	return claims.PhoneNumber, claims.Role, nil
}
