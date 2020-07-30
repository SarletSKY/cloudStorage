package handler

import (
	"filestore-server-study/common"
	"filestore-server-study/middleware"
	"filestore-server-study/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// token 拦截器
func HTTPInterceptor(c *gin.Context) {
	username := c.Request.FormValue("username")
	token := c.Request.FormValue("token")

	c.Abort() //报错后面的方法不用在执行
	// 验证token
	if len(username) < 3 || !middleware.ValidToToken(token) {
		resp := util.NewRespMsg(int(common.StatusInvalidToken), "token无效", nil)
		c.JSON(http.StatusOK, resp)
		return
	}
	c.Next()
}

// 允许跨域
func CORS(c *gin.Context) {
	c.Writer.Header().Set("text/plain", "*")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

	if c.Request.Method == "OPTIONS" {
		c.String(http.StatusOK, "")
	}

	// 调用下个中间件
	c.Next()
}
