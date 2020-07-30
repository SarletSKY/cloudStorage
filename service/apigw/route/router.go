package route

import (
	"filestore-server-study/service/apigw/handler"
	"github.com/gin-gonic/gin"

	"filestore-server-study/middleware"
)

// Router: 网关api路由
func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	// 注册
	router.GET("/user/signup", handler.SignUpUser)
	router.GET("/user/signin", handler.SignInUser)
	router.POST("/user/signup", handler.DoSignUpUser)
	router.POST("/user/signin", handler.DoSignInUser)

	router.Use(middleware.CORS)
	//中间件
	router.Use(middleware.HTTPInterceptor())

	router.POST("/user/info", handler.QueryUserInfo)
	router.POST("/file/query", handler.GetManyFileMetaInfo)
	router.POST("/file/update", handler.UpdateFileInfo)

	return router
}
