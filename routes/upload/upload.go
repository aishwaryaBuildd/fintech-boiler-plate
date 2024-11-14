package upload

import (
	"fintech/controllers/upload"
	"fintech/pkg/gcp"
	"fintech/pkg/vdo"
	"fintech/store"

	"github.com/gin-gonic/gin"
)

func UploadRoutes(r *gin.Engine, db store.Store, storageClient gcp.GCPUploader, vdo *vdo.VideoCipherClient) {
	controller := upload.Controller{Store: db, StorageClient: storageClient, VDO: vdo}

	r.POST("/upload/gcp/files", controller.UploadGCP)
	r.POST("/upload/vdo", controller.UploadVDO)
	r.POST("/view/vdo/:id", controller.ViewVDO)
	r.POST("/upload/vdo/:id/thumbnail", controller.VDOThumbnail)
	r.POST("/upload/vdo/:id/caption", controller.VDOCaption)
}
