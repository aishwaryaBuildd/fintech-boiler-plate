package server

import (
	"github.com/aishwaryaBuildd/go_boiler_plate/store"

	"github.com/gin-gonic/gin"
)

func NewServer(store store.Store) *gin.Engine {
	r := gin.Default()

	h := NewUserHandler(store)

	AttachUserRoutes(h, r)

	return r
}

func AttachUserRoutes(h *UserHandler, r *gin.Engine) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}
