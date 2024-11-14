package lectures

import (
	"fintech/controllers/lectures"
	"fintech/middlewares"
	"fintech/pkg/vdo"
	"fintech/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func lecturesRoutes(r *gin.Engine, db store.Store, VDO *vdo.VideoCipherClient) {
	controller := lectures.Controller{Store: db, VDO: VDO}

	r.POST("/courses/:id/sections/:section_id/lectures", middlewares.AdminMiddleware, courseMiddleware(db), controller.Create)
	r.GET("/courses/:id/sections/:section_id/lectures", middlewares.AuthMiddleware, courseMiddleware(db), controller.List)
	r.GET("/courses/:id/sections/:section_id/:section_id/lectures/:lecture_id", middlewares.AuthMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.Get)
	r.PATCH("/courses/:id/sections/:section_id/lectures/:lecture_id", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.Update)
	r.DELETE("/courses/:id/sections/:section_id/lectures/:lecture_id", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.Delete)

	r.POST("/courses/:id/sections/:section_id/lectures/:lecture_id/upload", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.Upload)
	r.POST("/courses/:id/sections/:section_id/lectures/:lecture_id/video", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.SetVideo)
	r.DELETE("/courses/:id/sections/:section_id/lectures/:lecture_id/video", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.RemoveVideo)
	r.POST("/courses/:id/sections/:section_id/lectures/:lecture_id/files", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.SetFiles)
	r.DELETE("/courses/:id/sections/:section_id/lectures/:lecture_id/files", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.RemoveFiles)
	r.POST("/courses/:id/sections/:section_id/lectures/:lecture_id/notes", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.SetNotes)
	r.DELETE("/courses/:id/sections/:section_id/lectures/:lecture_id/notes", middlewares.AdminMiddleware, courseMiddleware(db), sectionMiddleware(db), lectureMiddleware(db), controller.RemoveNotes)
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
		section, err := db.GetSection(c, section_id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.Set("section", section)
	}
}

func lectureMiddleware(db store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		lecture_id := c.Param("lecture_id")
		lecture, err := db.GetLecture(c, lecture_id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.Set("lecture", lecture)
	}
}
