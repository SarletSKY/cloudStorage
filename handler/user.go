package handler

import (
	"filestore-server-study/config"
	"filestore-server-study/db"
	"filestore-server-study/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// SignUpUser: 注册用户[GET]
func SignUpUser(c *gin.Context) {
	c.Redirect(http.StatusFound, "http://"+c.Request.Host+"/static/view/signup.html")
}

// DoSignUpUser: 注册用户[POST]
func DoSignUpUser(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	// 判断用户密码的正确性
	if len(username) < 3 || len(password) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Invalid parameter",
		})
		return
	}
	// 对密码进行进密
	encPassword := util.Sha1([]byte(password + config.PasswordSalt))

	// 将数据加到数据库
	suc := db.SignUpUserdb(username, encPassword)
	fmt.Println(suc)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"msg":     "注册成功",
			"data":    nil,
			"forward": "/user/signin",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  "注册失败",
		})
	}
}

// SignInUser: 登录用户[GET]
func SignInUser(c *gin.Context) {
	c.Redirect(http.StatusFound, "http://"+c.Request.Host+"/static/view/signin.html")
}

// DoSignInUser: 登录用户[POST]Salt
func DoSignInUser(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	encPwd := util.Sha1([]byte(password + config.PasswordSalt))
	suc := db.SignInUserdb(username, encPwd)
	if !suc {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "密码校验失败",
			"data": nil,
		})
		return
	}

	// 生成token凭证 [token 40位  md5加密]
	token := GetToken(username)
	suc = db.RegisterTokendb(username, token)
	if !suc {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "登录失败",
			"data": nil,
		})
		return
	}
	// 重定向到页面
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + c.Request.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "octet-stream", resp.JSONBytes())
}

// 获取token
func GetToken(username string) string {
	timeNow := fmt.Sprintf("%x", time.Now().Unix())
	tokenSalt := "private_zwx_key"
	encToken := util.MD5([]byte(username + timeNow + tokenSalt))
	return encToken + timeNow[:8]
}

//查询用户信息 [这里要返回到前端两个数据,username与注册时间]
func QueryUserInfo(c *gin.Context) {

	username := c.Request.FormValue("username")
	/*token := request.Form.Get("token") // ????????打印下
	// 检验token
	suc := ValidToToken(token)
	if !suc {
		writer.WriteHeader(http.StatusForbidden)
		return
	}*/
	// 向数据库查询信息
	user, err := db.QueryUserInfodb(username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}
	// 返回信息到前端
	resp := util.RespMsg{
		Code: http.StatusOK,
		Msg:  "OK",
		Data: user,
	}
	c.Data(http.StatusOK, "octet-stream", resp.JSONBytes())
}
