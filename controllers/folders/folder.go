package folders

import (
	"fintech/pkg/vdo"
	"fintech/store"
	"fintech/store/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Controller struct {
	Store store.Store
	VDO   *vdo.VideoCipherClient
}

func (controller Controller) Create(c *gin.Context) {
	var req mutateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	course := c.MustGet("course").(models.Course)

	newUUID := uuid.New()
	vdoFolder, err := controller.VDO.CreateSubFolder(newUUID.String(), course.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	folder := models.Folder{
		ID:          newUUID,
		Name:        req.Name,
		Description: req.Description,
		FolderID:    vdoFolder.ID,
		CourseID:    course.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = controller.Store.CreateFolder(c, folder)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Folder ID already exists"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusCreated, folder)
}

func (controller Controller) Update(c *gin.Context) {
	folder := c.MustGet("folder").(models.Folder)
	var req mutateRequest
	req.Description = folder.Description
	req.Name = folder.Name
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	folder.Description = req.Description
	folder.Name = req.Name
	err := controller.Store.UpdateFolder(c, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, folder)
}

func (controller Controller) List(c *gin.Context) {
	folders, err := controller.Store.ListFolder(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, folders)
}

func (controller Controller) Get(c *gin.Context) {
	folder := c.MustGet("folder").(models.Folder)

	vdoFolder, err := controller.VDO.GetSubFolders(folder.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	resp := CourseDetailedResponse{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		Folder:      *vdoFolder,
		CourseID:    folder.CourseID.String(),
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)

}

func (controller Controller) Delete(c *gin.Context) {
	folder := c.MustGet("folder").(models.Folder)
	err := controller.Store.DeleteFolder(c, folder.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller Controller) Upload(c *gin.Context) {
	folder := c.MustGet("folder").(models.Folder)

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

	credentials, err := controller.VDO.GetUploadCredentials(videoTitle, folder.FolderID)
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

type mutateRequest struct {
	Name        string `json:"name" validate:"min=5,max=50"`
	Description string `json:"description" validate:"min=5,max=500"`
}

type CourseDetailedResponse struct {
	ID          uuid.UUID          `db:"id"`          // Matches CHAR(36) for UUID
	Name        string             `db:"name"`        // VARCHAR(50), non-nullable
	Description string             `db:"description"` // VARCHAR(300), nullable, use sql.NullString
	CourseID    string             `db:"course_id"`   // INT, non-nullable
	Folder      vdo.FolderResponse `db:"folder"`
	CreatedAt   time.Time          `db:"created_at"` // DATETIME(6), default CURRENT_TIMESTAMP(6)
	UpdatedAt   time.Time          `db:"updated_at"`
}
