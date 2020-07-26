package db

import (
	"filestore-server-study/db/mysql"
	"fmt"
)

// 数据来的用户表结构体
type TableUser struct {
	Username     string
	Phone        string
	Email        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// 向数据库注册用户
func SignUpUserdb(username string, password string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user (`user_name`,`user_pwd`) values(?,?)")

	if err != nil {
		fmt.Println("Failed to prepare statement register user data ,err", err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("sql exec failed,err: ", err.Error())
		return false
	}

	// 查询数据库是否已经存在该用户
	if rf, err := result.RowsAffected(); err == nil && rf >= 0 {
		return true
	}

	return false
}

// 向数据库登录用户
func SignInUserdb(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement signIn user data ,err", err.Error())
		return false
	}

	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("sql exec failed,err: ", err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found" + username)
		return false
	}

	// 判断密码
	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}

	return false
}

// 注册token
func RegisterTokendb(username string, token string) bool {
	// 这里使用replace
	stmt, err := mysql.DBConn().Prepare("replace into tbl_user_token (`user_name`,`user_token`) values(?,?)")

	if err != nil {
		fmt.Println("Failed to prepare statement register token data ,err", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("sql exec failed,err: ", err.Error())
		return false
	}
	return true
}

// 查询用户信息
func QueryUserInfodb(username string) (TableUser, error) {
	user := TableUser{}
	stmt, err := mysql.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement register query user data ,err", err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println("sql exec failed,err: ", err.Error())
		return user, err
	}
	return user, nil
}
