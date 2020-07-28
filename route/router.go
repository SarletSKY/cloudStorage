package route

import (
	"filestore-server-study/handler"
	"github.com/gin-gonic/gin"
)

// 新建gin框架的路由规则
func Router() *gin.Engine {

	// gin初始路由
	router := gin.Default()

	//加载静态资源
	router.Static("/static/", "./static")

	// 不需要验证的路由接口
	router.GET("/user/signup", handler.SignUpUser)
	router.GET("/user/signin", handler.SignInUser)
	router.POST("/user/signup", handler.DoSignUpUser)
	router.POST("/user/signin", handler.DoSignInUser)
	router.POST("/user/info", handler.QueryUserInfo)

	// 加载中间件
	//router.Use(http.HTTPInterceptor)
	router.Use(handler.CORS)
	router.Use(handler.HTTPInterceptor)

	// 需要验证的路由接口
	router.GET("/file/upload", handler.UploadHandler)
	router.GET("/file/meta", handler.GetFileMetaInfo)
	router.GET("/file/upload/success", handler.UploadSucHandler)
	router.POST("/file/upload", handler.DoUploadHandler)
	router.POST("/file/query", handler.GetManyFileMetaInfo)
	router.POST("/file/update", handler.UpdateFileInfo)
	router.POST("/file/delete", handler.DeleteFile)
	router.POST("/file/fastupload", handler.FastUploadUserFile)

	// 文件下载
	router.GET("/file/download", handler.DownLoadFile)
	router.POST("/file/download", handler.DownLoadFile)
	router.POST("/file/downloadurl", handler.DownloadURL)
	router.POST("/file/download/range", handler.RangeDownload)

	//分块下载
	router.POST("/file/mpupload/init", handler.InitMultipartUpload)
	router.POST("/file/mpupload/uppart", handler.MultipartUpload)
	router.POST("/file/mpupload/complete", handler.CompleteMultipartUpload)
	router.POST("/file/mpupload/delete", handler.CancelUpload)

	return router
}
