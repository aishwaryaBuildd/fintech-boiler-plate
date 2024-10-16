package routes

import (
	"database/sql"
	"fintech/controller"

	"github.com/gin-gonic/gin"
)

func UserActionRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/video-analytics/:videoId", func(c *gin.Context) {
		controller.GetVideoAnalytics(c, db)
	})
}
