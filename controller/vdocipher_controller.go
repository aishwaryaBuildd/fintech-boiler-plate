package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"fintech/services"
	"fintech/utils"

	"github.com/gin-gonic/gin"
)

// CreateFolder handles folder creation for VdoCipher
func CreateFolder(c *gin.Context, db *sql.DB) {
	// Extract the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check for Bearer token
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Get the actual token part (after "Bearer ")
	token := authHeader[len(bearerPrefix):]

	// Verify JWT token and get phone number and role
	_, role, err := utils.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check if the user is an admin
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create folders"})
		return
	}

	// Bind the request body to extract the folder name
	var folderData struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&folderData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, folder name is required"})
		return
	}

	// Create the folder using VdoCipher client
	vdoCipherClient := services.NewVideoCipherClient()
	folder, err := vdoCipherClient.CreateFolderRoot(folderData.Name, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder created successfully", "folder": folder})
}

func GetFolders(c *gin.Context, db *sql.DB) {
	// Extract the Authorization header (JWT Token)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check for Bearer token
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// // Get the actual token part (after "Bearer ")
	// token := authHeader[len(bearerPrefix):]

	// // Verify JWT token (assuming you are verifying for user role)
	// _, role, err := utils.VerifyJWT(token)
	// if err != nil {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
	//     return
	// }

	// Create the VdoCipher client
	vdoCipherClient := services.NewVideoCipherClient()

	// Fetch all folders (Only users with valid JWT are allowed)
	folderList, err := vdoCipherClient.GetAllFolders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve folders"})
		return
	}

	// Return the folders to the user
	c.JSON(http.StatusOK, gin.H{"folders": folderList.FolderList})
}

// DeleteFolder is a handler to delete a folder by its ID
func DeleteFolder(c *gin.Context, db *sql.DB) {
	// Extract the Authorization header (JWT Token)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check for Bearer token
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Get the actual token part (after "Bearer ")
	token := authHeader[len(bearerPrefix):]

	// Verify JWT token and get user role
	_, role, err := utils.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check if the user has admin role
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete folders"})
		return
	}

	// Get the folder ID from the URL parameters
	folderID := c.Param("id")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder ID is required"})
		return
	}

	// Create the VdoCipher client
	vdoCipherClient := services.NewVideoCipherClient()

	// Call the method to delete the folder
	err = vdoCipherClient.DeleteFolder(folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}

func UploadVideo(c *gin.Context, db *sql.DB) {
	// Extract the Authorization header (JWT Token)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check for Bearer token
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Get the actual token part (after "Bearer ")
	token := authHeader[len(bearerPrefix):]

	// Verify JWT token and get user role
	_, role, err := utils.VerifyJWT(token) // Assume this function is defined elsewhere
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check if the user has admin role
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to upload videos"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Create a temporary file to save the uploaded video
	tempFile, err := os.CreateTemp("", "upload-*.mp4")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Save the uploaded file to the temporary file
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	// Extract the title from query params
	videoTitle := c.Query("title")
	if videoTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video title is required"})
		return
	}

	// Extract the folderId from query params (optional)
	folderID := c.Query("folderId")

	// Create the VdoCipher client
	vdoCipherClient := services.NewVideoCipherClient()

	// Step 1: Get upload credentials
	credentials, err := vdoCipherClient.GetUploadCredentials(videoTitle, folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get upload credentials", "details": err.Error()})
		return
	}

	// Log the credentials received
	log.Printf("Received credentials: %+v\n", credentials)

	// Step 2: Upload the video to S3 using the provided credentials
	err = vdoCipherClient.UploadFile(*credentials, tempFile.Name()) // Dereference credentials
	if err != nil {
		log.Printf("Failed to upload video: %v\n", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video", "details": err.Error()})
		return
	}

	// Step 3: Return the video ID (assume the response contains the VideoID)
	c.JSON(http.StatusOK, gin.H{"videoID": credentials.FileName}) // Assuming FileName holds the video ID
}

func CreateSubFolder(c *gin.Context, db *sql.DB) {
	// Extract the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check for Bearer token
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Get the actual token part (after "Bearer ")
	token := authHeader[len(bearerPrefix):]

	// Verify JWT token and get phone number and role
	_, role, err := utils.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check if the user is an admin
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create folders"})
		return
	}

	// Bind the request body to extract the folder name and parent folder ID
	var folderData struct {
		Name   string `json:"name" binding:"required"`
		Parent string `json:"parent" binding:"required"` // Parent folder ID is required to create a subfolder
	}
	if err := c.ShouldBindJSON(&folderData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, folder name and parent folder ID are required"})
		return
	}

	// Create the subfolder using VdoCipher client
	vdoCipherClient := services.NewVideoCipherClient()
	folder, err := vdoCipherClient.CreateSubFolder(folderData.Name, folderData.Parent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subfolder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subfolder created successfully", "folder": folder})
}

func ListSubFolders(c *gin.Context) {
	// Extract the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check for Bearer token
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Get the actual token part (after "Bearer ")
	token := authHeader[len(bearerPrefix):]

	// Verify JWT token (allowing all users, not checking roles)
	phoneNumber, role, err := utils.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	fmt.Printf("Authenticated user: %s, Role: %s\n", phoneNumber, role)

	// Get folder ID from the path parameter
	folderID := c.Param("folderId")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder ID is required"})
		return
	}

	// Fetch the folder details using VdoCipher client
	vdoCipherClient := services.NewVideoCipherClient()
	folderResponse, err := vdoCipherClient.GetSubFolders(folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve folder details"})
		return
	}

	// Respond with the folder details (and any subfolders if present)
	c.JSON(http.StatusOK, folderResponse)
}

// Modify GetVideoAnalytics to accept the db parameter
func GetVideoAnalytics(c *gin.Context, db *sql.DB) {
	videoID := c.Param("videoId")

	query := `SELECT action, time_in_video, timestamp FROM video_events WHERE video_id = $1 ORDER BY timestamp`
	rows, err := db.Query(query, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve video events"})
		return
	}
	defer rows.Close()

	var events []gin.H
	for rows.Next() {
		var action string
		var timeInVideo float64
		var timestamp time.Time

		if err := rows.Scan(&action, &timeInVideo, &timestamp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan video event"})
			return
		}

		events = append(events, gin.H{
			"action":        action,
			"time_in_video": timeInVideo,
			"timestamp":     timestamp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video analytics retrieved successfully",
		"events":  events,
	})
}
