package upload

import (
	"context"
	"fintech/pkg/vdo"
	"fintech/store"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Store         store.Store
	StorageClient CloudUploader
	VDO           vdo.VideoCipherClient
}

func (controller *Controller) UploadGCP(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	defer file.Close()

	url, err := controller.StorageClient.UploadFile(c, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (controller Controller) UploadVDO(c *gin.Context) {
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

	videoTitle := c.Query("title")
	if videoTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video title is required"})
		return
	}

	credentials, err := controller.VDO.GetUploadCredentials(videoTitle, "root")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Step 2: Upload the video to S3 using the provided credentials
	err = controller.VDO.UploadFile(*credentials, tempFile.Name()) // Dereference credentials
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videoID": credentials.FileName})
}

func (controller *Controller) ViewVDO(c *gin.Context) {
	videoID := c.Param("id")
	otpResponse, err := controller.VDO.GenerateVideoOTP(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": otpResponse.PlayURL})

}

type CloudUploader interface {
	UploadFile(context context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
