package sections

import (
	sectionController "fintech/controllers/sections"
	"fintech/middlewares"
	"fintech/pkg/vdo"
	"fintech/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SectionRoutes(r *gin.Engine, db store.Store, VDO *vdo.VideoCipherClient) {
	controller := sectionController.Controller{Store: db, VDO: VDO}

	r.POST("/courses/:id/sections", middlewares.AdminMiddleware, courseMiddleware(db), controller.Create)
	r.GET("/courses/:id/sections", middlewares.AuthMiddleware, courseMiddleware(db), controller.List)
	r.GET("/courses/:id/sections/:section_id", middlewares.AuthMiddleware, courseMiddleware(db), sectionMiddleware(db), controller.Get)
	r.PATCH("/courses/:id/sections/:section_id", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), controller.Update)
	r.DELETE("/courses/:id/sections/:section_id", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), controller.Delete)

}

func courseMiddleware(db store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		course, err := db.GetCourse(c, courseID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.Set("course", course)
	}

}

func sectionMiddleware(db store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		section_id := c.Param("section_id")
		folder, err := db.GetSection(c, section_id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.Set("section", folder)
	}

}
