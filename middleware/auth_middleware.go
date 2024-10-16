// middleware/auth.go

package middleware

import (
	"net/http"

	"fintech/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		phoneNumber, _, err := utils.VerifyJWT(token) // We only care about phone number here
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("phone_number", phoneNumber) // Store phone number in context
		c.Next()
	}
}

// AdminMiddleware checks if the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assume we have a way to determine if the user is an admin
		isAdmin := true // Placeholder for actual admin check logic
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
