package routes

import (
	"database/sql"
	"fintech/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/register", func(c *gin.Context) { controller.Register(c, db) })
	r.POST("/verify", func(c *gin.Context) { controller.Verify(c, db) })
}
