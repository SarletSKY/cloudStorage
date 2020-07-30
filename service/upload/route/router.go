package route

import (
	"filestore-server-study/middleware"
	"filestore-server-study/service/upload/api"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	router.Use(middleware.CORS)

	router.POST("/file/upload", api.DoUploadHandler)
	router.POST("/file/fastupload", api.FastUploadUserFile)

	//分块下载
	router.POST("/file/mpupload/init", api.InitMultipartUpload)
	router.POST("/file/mpupload/uppart", api.MultipartUpload)
	router.POST("/file/mpupload/complete", api.CompleteMultipartUpload)
	//router.POST("/file/mpupload/delete", handler.CancelUpload)
	return router
}
