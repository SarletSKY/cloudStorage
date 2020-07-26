package handler

import (
	"filestore-server-study/db"
	"filestore-server-study/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//设置盐
const pwdSalt = "zwx"

// 注册用户
func SignUpUser(writer http.ResponseWriter, request *http.Request) {
	// 如果是get请求，加载页面
	if request.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 读出页面
		writer.Write(data)
		return
	}

	// 解析参数
	request.ParseForm()

	username := request.Form.Get("username")
	password := request.Form.Get("password")

	// 判断用户密码的正确性
	if len(username) < 3 || len(password) < 5 {
		writer.Write([]byte("Invalid parameter"))
		return
	}

	// 对密码进行进密
	encPassword := util.Sha1([]byte(password + pwdSalt))

	// 将数据加到数据库
	suc := db.SignUpUserdb(username, encPassword)
	fmt.Println(suc)
	if suc {
		writer.Write([]byte("SUCCESS"))
	} else {
		writer.Write([]byte("FAILED"))
	}
}

// 登录用户
func SignInUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		http.Redirect(writer, request, "/static/view/signin.html", http.StatusFound)
		return
	}

	// 登录用户名/校验密码
	request.ParseForm()

	username := request.Form.Get("username")
	password := request.Form.Get("password")
	encPwd := util.Sha1([]byte(password + pwdSalt))
	suc := db.SignInUserdb(username, encPwd)
	if !suc {
		writer.Write([]byte("FAILED"))
		return
	}

	// 生成token凭证 [token 40位  md5加密]
	token := GetToken(username)
	suc = db.RegisterTokendb(username, token)
	if !suc {
		writer.Write([]byte("FAILED"))
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
			Location: "http://" + request.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	writer.Write(resp.JSONBytes())
}

// 获取token
func GetToken(username string) string {
	timeNow := fmt.Sprintf("%x", time.Now().Unix())
	tokenSalt := "private_zwx_key"
	encToken := util.MD5([]byte(username + timeNow + tokenSalt))
	return encToken + timeNow[:8]
}

//查询用户信息 [这里要返回到前端两个数据,username与注册时间]
func QueryUserInfo(writer http.ResponseWriter, request *http.Request) {
	// 解析请求参数
	request.ParseForm()

	username := request.Form.Get("username")
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
		writer.WriteHeader(http.StatusForbidden)
		return
	}
	// 返回信息到前端
	resp := util.RespMsg{
		Code: 200,
		Msg:  "OK",
		Data: user,
	}
	writer.Write(resp.JSONBytes())
}

// 验证token
func ValidToToken(token string) bool {
	//token是否为40位
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// example，假设token的有效期为1天   (根据同学们反馈完善, 相对于视频有所更新)
	tokenTS := token[:8]
	if util.Hex2Dec(tokenTS) < time.Now().Unix()-86400 {
		return false
	}

	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致

	return true
}
