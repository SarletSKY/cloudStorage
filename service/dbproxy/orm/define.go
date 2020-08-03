package orm

import "database/sql"

// 对应数据库的表结构体
// 表的文件结构体
type TableFile struct {
	FileSha1 string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAdd  sql.NullString
}

// 用户表结构体
type TableUser struct {
	Username     string
	Phone        string
	Email        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// 用户文件表
type TableUserFile struct {
	UserName    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
	FileHash    string
}

// sql执行返回结构体
type ExecResult struct {
	Suc  bool        `json:"suc"`  // 成功与否
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  //错误信息
	Data interface{} `json:"data"` // 返回数据信息
}
