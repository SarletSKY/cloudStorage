package handler

import (
	"context"
	"filestore-server-study/common"
	"filestore-server-study/config"
	userProto "filestore-server-study/service/account/proto"
	"filestore-server-study/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"log"
	"net/http"
)

var (
	userCli userProto.UserService // 全局的user服务
)

func init() {
	// micro获取一个服务
	service := micro.NewService(micro.Registry(config.RegistryConsul()))
	// 初始化micro客户端
	service.Init()

	// 初始化account客户端
	userCli = userProto.NewUserService("go.micro.service.user", service.Client())
}

// 注册与登录服务
// SignUpUser: 注册用户[GET]
func SignUpUser(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DoSignUpUser: 注册用户[POST]
func DoSignUpUser(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	resp, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})
}

// SignInUser: 登录用户[GET]
func SignInUser(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// DoSignInUser: 登录用户[POST]Salt
func DoSignInUser(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	// 用户登录
	userCli, err := userCli.Signin(context.TODO(), &userProto.ReqSignin{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if userCli.Code != common.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"code": userCli.Code,
			"msg":  "登录失败",
		})
		return
	}

	cliResp := util.RespMsg{
		Code: int(common.StatusOK),
		Msg:  "登录成功",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    userCli.Token,
		},
	}
	c.Data(http.StatusOK, "octet-stream", cliResp.JSONBytes())
}

//查询用户信息 [这里要返回到前端两个数据,username与注册时间]
func QueryUserInfo(c *gin.Context) {
	fmt.Println("?")
	username := c.Request.FormValue("username")
	// 向数据库查询信息
	resp, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: username,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 返回信息到前端
	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: gin.H{
			"Username": username,
			"SignupAt": resp.SignupAt,
			// TODO: 完善其他字段信息
			"LastActive": resp.LastActiveAt,
		},
	}
	c.Data(http.StatusOK, "octet-stream", cliResp.JSONBytes())
}
