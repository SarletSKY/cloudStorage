package middleware

import (
	"filestore-server-study/common"
	"filestore-server-study/db"
	"filestore-server-study/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// token 拦截器
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		// 验证token
		if len(username) < 3 || !ValidToToken(username, token) {
			c.Abort() //报错后面的方法不用在执行
			resp := util.NewRespMsg(int(common.StatusInvalidToken), "token无效", nil)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
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

// 验证token
func ValidToToken(username string, token string) bool {
	//token是否为40位
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// example，假设token的有效期为1天   (根据同学们反馈完善, 相对于视频有所更新)
	tokenTS := token[32:]
	if util.Hex2Dec(tokenTS) < time.Now().Unix()-86400 {
		log.Println("token expired: " + token)
		return false
	}

	tokenToDB, err := db.GetUserToken(username)
	if err != nil || tokenToDB != token {
		return false
	}
	return true
}
