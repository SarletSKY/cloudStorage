package orm

import (
	_ "database/sql"
	"filestore-server-study/service/dbproxy/mysql"
	"log"
)

// 向数据库注册用户
func SignUpUserdb(username string, password string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user (`user_name`,`user_pwd`) values(?,?)")

	if err != nil {
		log.Println("Failed to prepare statement register user data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	// 查询数据库是否已经存在该用户
	if rf, err := result.RowsAffected(); err == nil && rf >= 0 {
		res.Suc = true
		return
	}
	res.Suc = false
	res.Msg = "用户已经存在"
	return
}

// 向数据库登录用户
func SignInUserdb(username string, encpwd string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println("Failed to prepare statement signIn user data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil {
		log.Println("username not found" + username)
		res.Suc = false
		res.Msg = "用户未注册"
		return
	}

	// 判断密码
	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		res.Suc = true
		res.Data = true
		return
	}

	res.Suc = false
	res.Msg = "用户名/密码不匹配"
	return
}

// 注册token
func RegisterTokendb(username string, token string) (res ExecResult) {
	// 这里使用replace
	stmt, err := mysql.DBConn().Prepare("replace into tbl_user_token (`user_name`,`user_token`) values(?,?)")

	if err != nil {
		log.Println("Failed to prepare statement register token data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// 查询用户信息
func QueryUserInfodb(username string) (res ExecResult) {
	user := TableUser{}
	stmt, err := mysql.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println("Failed to prepare statement register query user data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = user
	return
}

// 查询用户token信息
func GetUserToken(username string) (res ExecResult) {
	var token string
	stmt, err := mysql.DBConn().Prepare("select user_token from tbl_user_token where user_name=?")
	if err != nil {
		log.Println("Failed to prepare statement query token failed ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&token)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = map[string]string{
		"token": token,
	}
	return
}

// 查看有没有该用户
func UserExist(username string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select 1 from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = map[string]bool{
		"exists": rows.Next(),
	}
	return
}
