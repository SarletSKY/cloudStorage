package route

import (
	microHandler "filestore-server-study/service/apigw/handler"
	"github.com/gin-gonic/gin"

	"filestore-server-study/handler"
)

// Router: 网关api路由
func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	// 注册
	router.GET("/user/signup", microHandler.SignUpUser)
	router.GET("/user/signin", microHandler.SignInUser)
	router.POST("/user/signup", microHandler.DoSignUpUser)
	router.POST("/user/signin", microHandler.DoSignInUser)
	router.POST("/user/info", microHandler.QueryUserInfo)

	//中间件
	router.Use(handler.HTTPInterceptor)
	router.Use(handler.CORS)

	router.POST("/file/query", handler.GetManyFileMetaInfo)
	router.POST("/file/update", handler.UpdateFileInfo)

	return router
}
