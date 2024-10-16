package routes

import (
	"database/sql"
	"fintech/controller"

	"github.com/gin-gonic/gin"
)

func VideoRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/create-folder", func(c *gin.Context) { controller.CreateFolder(c, db) })
	r.GET("/get-folders", func(c *gin.Context) { controller.GetFolders(c, db) })
	r.DELETE("/delete-folder/:id", func(c *gin.Context) { controller.DeleteFolder(c, db) })
	r.POST("/upload-video", func(c *gin.Context) { controller.UploadVideo(c, db) })
	r.POST("/create-subfolder", func(c *gin.Context) { controller.CreateSubFolder(c, db) })
	r.GET("/list-subfolders/:folderId", func(c *gin.Context) { controller.ListSubFolders(c) })

}
