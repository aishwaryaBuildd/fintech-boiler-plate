package upload

import (
	"context"
	"fintech/store"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Store         store.Store
	StorageClient CloudUploader
}

func (controller *Controller) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	defer file.Close()

	url, err := controller.StorageClient.UploadImage(c, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})

}

type CloudUploader interface {
	UploadImage(context context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
