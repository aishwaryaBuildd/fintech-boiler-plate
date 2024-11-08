package folders

import (
	folderController "fintech/controllers/folders"
	"fintech/middlewares"
	"fintech/pkg/vdo"
	"fintech/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FolderRoutes(r *gin.Engine, db store.Store, VDO *vdo.VideoCipherClient) {
	controller := folderController.Controller{Store: db, VDO: VDO}

	r.POST("/courses/:id/folders", middlewares.AdminMiddleware, courseMiddleware(db), controller.Create)
	r.GET("/courses/:id/folders", middlewares.AuthMiddleware, courseMiddleware(db), controller.List)
	r.GET("/courses/:id/folders/:folder_id", middlewares.AuthMiddleware, courseMiddleware(db), folderMiddleware(db), controller.Get)
	r.PATCH("/courses/:id/folders/:folder_id", middlewares.AdminMiddleware, courseMiddleware(db), folderMiddleware(db), controller.Update)
	r.DELETE("/courses/:id/folders/:folder_id", middlewares.AdminMiddleware, courseMiddleware(db), folderMiddleware(db), controller.Delete)

	r.POST("/courses/:id/folders/:folder_id/upload", middlewares.AdminMiddleware, courseMiddleware(db), folderMiddleware(db), controller.Upload)

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

func folderMiddleware(db store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		folder_id := c.Param("folder_id")
		folder, err := db.GetFolder(c, folder_id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.Set("folder", folder)
	}

}
