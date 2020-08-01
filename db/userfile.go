package db

import (
	"filestore-server-study/db/mysql"
	"fmt"
)

type TableUserFile struct {
	UserName    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
	FileHash    string
}

// 上传文件到user_file表
func OnUserFileUploadFinshedDB(username, filename, filehash string, filesize int64, uploadAt string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`,`file_name`,`file_size`,`file_sha1`,`upload_at`,`status`) " +
		" values (?,?,?,?,?,1)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filename, filesize, filehash, uploadAt)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// 批量查询用户文件接口
func QueryMantUserFileDB(username string, limit int) ([]TableUserFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_name,file_sha1,file_size,upload_at,last_update from tbl_user_file where user_name=? and status!=2 limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	userFileArr := make([]TableUserFile, 0)
	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	// 返回多条数据
	for rows.Next() {
		tuf := TableUserFile{}
		err = rows.Scan(&tuf.FileName, &tuf.FileHash, &tuf.FileSize, &tuf.UploadAt, &tuf.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		userFileArr = append(userFileArr, tuf)
	}
	return userFileArr, nil
}

// 修改用户文件表的元信息
func UpdateUserFileInfoDB(username string, filename string, filehash string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_user_file set file_name=? where user_name=? and file_sha1=? limit 1")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// 判断文件是否存在
func IsUserFileUpload(username string, filehash string) bool {
	stmt, err := mysql.DBConn().Prepare("select 1 from tbl_user_file where user_name=? and file_sha1=? and status=1 limit 1")
	if err != nil {
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		return false
	} else if rows == nil || !rows.Next() {
		return false
	}
	return true
}

// 查询单个用户文件元信息
func QueryUserFileDB(username string, filehash string) (*TableUserFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_name,file_sha1,file_size,upload_at,last_update from tbl_user_file where user_name=? and file_sha1=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	ufile := TableUserFile{}
	for rows.Next() {
		err = rows.Scan(&ufile.FileName, &ufile.FileHash, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			return nil, err
		}
	}
	return &ufile, nil
}

// 删除用户文件表
func DeleteUserFileDB(username string, filehash string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_user_file set status=2 where user_name=? and file_sha1=? limit 1")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// 用户文件名改之前，文件名有没有被使用
func QueryUserFileNameExist(username string, filename string) bool {
	var queryFileName string
	stmt, err := mysql.DBConn().Prepare("select file_name from tbl_user_file where user_name=? and file_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()
	err = stmt.QueryRow(username, filename).Scan(&queryFileName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if queryFileName != "" {
		return true
	}
	return false
}
