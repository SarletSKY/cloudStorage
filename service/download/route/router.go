package route

import (
	"filestore-server-study/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/", "./static")

	router.Use(handler.CORS)

	router.GET("/file/download", handler.DownLoadFile)
	router.GET("/file/download/range", handler.RangeDownload)
	router.POST("/file/downloadurl", handler.DownloadURL)

	return router
}
