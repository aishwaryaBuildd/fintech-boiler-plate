package courses

import (
	courseController "fintech/controllers/courses"
	"fintech/middlewares"
	"fintech/pkg/vdo"
	"fintech/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(r *gin.Engine, db store.Store, VDO *vdo.VideoCipherClient) {
	controller := courseController.Controller{Store: db, VDO: VDO}

	r.POST("/courses", middlewares.AdminMiddleware, controller.Create)
	r.GET("/courses", middlewares.AuthMiddleware, controller.List)
	r.GET("/courses/:id", middlewares.AuthMiddleware, courseMiddleware(db), controller.Get)
	r.PATCH("/courses/:id", middlewares.AdminMiddleware, courseMiddleware(db), controller.Update)
	r.DELETE("/courses/:id", middlewares.AdminMiddleware, courseMiddleware(db), controller.Delete)
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
