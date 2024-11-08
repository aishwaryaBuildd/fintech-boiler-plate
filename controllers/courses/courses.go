package courses

import (
	"fintech/pkg/vdo"
	"fintech/store"
	"fintech/store/models"
	"net/http"
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

	newUUID := uuid.New()
	vdoFolder, err := controller.VDO.CreateFolderRoot(newUUID.String(), "root")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	course := models.Course{
		ID:          newUUID,
		Name:        req.Name,
		Description: req.Description,
		FolderID:    vdoFolder.ID,
		AuthorID:    c.MustGet("user_id").(int),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = controller.Store.CreateCourse(c, course)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID already exists"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusCreated, course)
}

func (controller Controller) Update(c *gin.Context) {
	course := c.MustGet("course").(models.Course)
	var req mutateRequest
	req.Description = course.Description
	req.Name = course.Name
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	course.Description = req.Description
	course.Name = req.Name
	err := controller.Store.UpdateCourse(c, course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (controller Controller) List(c *gin.Context) {
	courses, err := controller.Store.ListCourse(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (controller Controller) Get(c *gin.Context) {
	course := c.MustGet("course").(models.Course)

	vdoFolder, err := controller.VDO.GetSubFolders(course.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	resp := CourseDetailedResponse{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		Folder:      *vdoFolder,
		AuthorID:    course.AuthorID,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

func (controller Controller) Delete(c *gin.Context) {
	course := c.MustGet("course").(models.Course)
	err := controller.Store.DeleteCourse(c, course.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusNoContent)
}

type mutateRequest struct {
	Name        string `json:"name" validate:"min=5,max=50"`
	Description string `json:"description" validate:"min=5,max=500"`
}

type CourseDetailedResponse struct {
	ID          uuid.UUID          `db:"id"`          // Matches CHAR(36) for UUID
	Name        string             `db:"name"`        // VARCHAR(50), non-nullable
	Description string             `db:"description"` // VARCHAR(300), nullable, use sql.NullString
	AuthorID    int                `db:"author_id"`   // INT, non-nullable
	Folder      vdo.FolderResponse `db:"folder"`
	CreatedAt   time.Time          `db:"created_at"` // DATETIME(6), default CURRENT_TIMESTAMP(6)
	UpdatedAt   time.Time          `db:"updated_at"`
}
