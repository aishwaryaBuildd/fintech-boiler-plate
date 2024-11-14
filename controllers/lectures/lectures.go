package lectures

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

	section := c.MustGet("section").(models.Section)

	newUUID := uuid.New()

	lecture := models.Lecture{
		ID:          newUUID,
		Name:        req.Name,
		Description: req.Description,
		SectionID:   section.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := controller.Store.CreateLecture(c, lecture)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Lecture ID already exists"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusCreated, section)
}

func (controller Controller) Update(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)
	var req mutateRequest
	req.Description = lecture.Description
	req.Name = lecture.Name
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	lecture.Description = req.Description
	lecture.Name = req.Name
	err := controller.Store.UpdateLecture(c, lecture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

func (controller Controller) List(c *gin.Context) {
	sections, err := controller.Store.ListLecture(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, sections)
}

func (controller Controller) Get(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)

	c.JSON(http.StatusOK, lecture)

}

func (controller Controller) Delete(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)
	err := controller.Store.DeleteLecture(c, lecture.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller Controller) SetVideo(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)

	var req videoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	lecture.VideoID = req.VideoID
	lecture.UpdatedAt = time.Now()

	err := controller.Store.UpdateLecture(c, lecture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

func (controller Controller) RemoveVideo(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)

	lecture.VideoID = ""
	lecture.UpdatedAt = time.Now()

	err := controller.Store.UpdateLecture(c, lecture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

func (controller Controller) SetFiles(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)

	var req filesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	lecture.FileName = req.FileName
	lecture.FileURL = req.URL
	lecture.UpdatedAt = time.Now()

	err := controller.Store.UpdateLecture(c, lecture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

func (controller Controller) RemoveFiles(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)

	lecture.FileName = ""
	lecture.FileURL = ""
	lecture.UpdatedAt = time.Now()

	err := controller.Store.UpdateLecture(c, lecture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

func (controller Controller) SetNotes(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)

	var req NotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	lecture.LectureNotes = req.Notes
	lecture.UpdatedAt = time.Now()

	err := controller.Store.UpdateLecture(c, lecture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

func (controller Controller) RemoveNotes(c *gin.Context) {
	lecture := c.MustGet("lecture").(models.Lecture)

	lecture.LectureNotes = ""
	lecture.UpdatedAt = time.Now()

	err := controller.Store.UpdateLecture(c, lecture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

type NotesRequest struct {
	// TODO: validate
	Notes string `json:"notes"`
}

type filesRequest struct {
	// TODO: validate
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

type videoRequest struct {
	// TODO: validate
	VideoID string `json:"video_id"`
}

func (controller Controller) Upload(c *gin.Context) {
	course := c.MustGet("course").(models.Course)

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

	credentials, err := controller.VDO.GetUploadCredentials(videoTitle, course.FolderID)
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
