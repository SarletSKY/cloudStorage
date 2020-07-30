package route

import (
	"filestore-server-study/middleware"
	"filestore-server-study/service/download/api"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/", "./static")

	router.Use(middleware.CORS)

	router.GET("/file/download", api.DownLoadFile)
	router.GET("/file/download/range", api.RangeDownload)
	router.POST("/file/downloadurl", api.DownloadURL)

	return router
}
