package route

import (
	"filestore-server-study/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	router.Use(handler.CORS)

	router.POST("/file/upload", handler.DoUploadHandler)
	router.POST("/file/fastupload", handler.FastUploadUserFile)

	//分块下载
	router.POST("/file/mpupload/init", handler.InitMultipartUpload)
	router.POST("/file/mpupload/uppart", handler.MultipartUpload)
	router.POST("/file/mpupload/complete", handler.CompleteMultipartUpload)
	router.POST("/file/mpupload/delete", handler.CancelUpload)
}
