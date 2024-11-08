package middlewares

import (
	"fintech/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid JWT token
func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	claims, err := utils.VerifyJWT(token) // We only care about phone number here
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("user_id", claims.UserID)           // Store phone number in context
	c.Set("phone_number", claims.PhoneNumber) // Store phone number in context
	c.Next()
}

// AdminMiddleware checks if the user is an admin
func AdminMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	claims, err := utils.VerifyJWT(token) // We only care about phone number here
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access"})
		c.Abort()
		return
	}
	c.Set("user_id", claims.UserID)           // Store phone number in context
	c.Set("phone_number", claims.PhoneNumber) // Store phone number in context
	c.Next()
}
