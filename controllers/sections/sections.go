package sections

import (
	"fintech/pkg/vdo"
	"fintech/store"
	"fintech/store/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	section := models.Section{
		ID:          newUUID,
		Name:        req.Name,
		Description: req.Description,
		CourseID:    course.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	c.JSON(http.StatusCreated, section)
}

func (controller Controller) Update(c *gin.Context) {
	section := c.MustGet("section").(models.Section)
	var req mutateRequest
	req.Description = section.Description
	req.Name = section.Name
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	section.Description = req.Description
	section.Name = req.Name
	err := controller.Store.UpdateSection(c, section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, section)
}

func (controller Controller) List(c *gin.Context) {
	sections, err := controller.Store.ListSection(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, sections)
}

func (controller Controller) Get(c *gin.Context) {
	section := c.MustGet("section").(models.Section)

	resp := CourseDetailedResponse{
		ID:          section.ID,
		Name:        section.Name,
		Description: section.Description,
		CourseID:    section.CourseID.String(),
		CreatedAt:   section.CreatedAt,
		UpdatedAt:   section.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)

}

func (controller Controller) Delete(c *gin.Context) {
	section := c.MustGet("section").(models.Section)
	err := controller.Store.DeleteSection(c, section.ID.String())
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
	ID          uuid.UUID `db:"id"`          // Matches CHAR(36) for UUID
	Name        string    `db:"name"`        // VARCHAR(50), non-nullable
	Description string    `db:"description"` // VARCHAR(300), nullable, use sql.NullString
	CourseID    string    `db:"course_id"`   // INT, non-nullable
	CreatedAt   time.Time `db:"created_at"`  // DATETIME(6), default CURRENT_TIMESTAMP(6)
	UpdatedAt   time.Time `db:"updated_at"`
}
