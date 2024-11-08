package auth

import (
	authController "fintech/controllers/auth"
	"fintech/store"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, db store.Store) {
	controller := authController.Controller{Store: db}

	r.POST("/register", controller.Register)
	r.POST("/verify", controller.Verify)
}
