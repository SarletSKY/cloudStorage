package handler

import (
	"context"
	"filestore-server-study/common"
	"filestore-server-study/config"
	"filestore-server-study/db"
	proto "filestore-server-study/service/account/proto"
	"filestore-server-study/util"
	"fmt"
	"time"
)

type User struct{}

// 获取token
func GetToken(username string) string {
	timeNow := fmt.Sprintf("%x", time.Now().Unix())
	tokenSalt := "private_zwx_key"
	encToken := util.MD5([]byte(username + timeNow + tokenSalt))
	return encToken + timeNow[:8]
}

// 注册用户
func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSignup) error {
	username := req.Username
	password := req.Password

	// 判断用户密码的正确性
	if len(username) < 3 || len(password) < 5 {
		resp.Code = common.StatusParamInvalid
		resp.Message = "注册参数无效"
		return nil
	}
	// 对密码进行进密
	encPassword := util.Sha1([]byte(password + config.PasswordSalt))

	// 将数据加到数据库
	suc := db.SignUpUserdb(username, encPassword)
	fmt.Println(suc)
	if suc {
		resp.Code = common.StatusOK
		resp.Message = "注册成功"
	} else {
		resp.Code = common.StatusRegisterFailed
		resp.Message = "注册失败"
	}
	return nil
}

// 登录用户
func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, resp *proto.RespSignin) error {
	username := req.Username
	password := req.Password
	encPwd := util.Sha1([]byte(password + config.PasswordSalt))
	suc := db.SignInUserdb(username, encPwd)
	if !suc {
		resp.Code = common.StatusLoginFailed
		return nil
	}

	// 生成token凭证 [token 40位  md5加密]
	token := GetToken(username)
	suc = db.RegisterTokendb(username, token)
	if !suc {
		resp.Code = common.StatusServerError
		return nil
	}
	resp.Code = common.StatusOK
	resp.Token = token
	return nil
}

// 查询用户信息
func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {
	username := req.Username
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
		resp.Code = common.StatusServerError
		resp.Message = "服务错误"
		return nil
	}

	if user.Username == "" {
		resp.Code = common.StatusUserNotExists
		resp.Message = "用户不存在"
		return nil
	}

	// 组装信息返回
	resp.Code = common.StatusOK
	resp.Username = user.Username
	resp.SignupAt = user.SignupAt
	resp.LastActiveAt = user.LastActiveAt
	resp.Status = int32(user.Status)
	resp.Email = user.Email
	resp.Phone = user.Phone
	// TODO: 需增加接口支持完善用户信息(email/phone等)
	return nil
}
