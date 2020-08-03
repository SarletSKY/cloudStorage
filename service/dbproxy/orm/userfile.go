package orm

import (
	_ "database/sql"

	"filestore-server-study/service/dbproxy/mysql"
	"log"
)

// 上传文件到user_file表
func OnUserFileUploadFinshedDB(username, filename, filehash string, filesize int64, uploadAt string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`,`file_name`,`file_size`,`file_sha1`,`upload_at`,`status`) " +
		" values (?,?,?,?,?,1)")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filename, filesize, filehash, uploadAt)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// 批量查询用户文件接口
func QueryMantUserFileDB(username string, limit int64) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select file_name,file_sha1,file_size,upload_at,last_update from tbl_user_file where user_name=? and status!=2 limit ?")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	userFileArr := make([]TableUserFile, 0)
	rows, err := stmt.Query(username, limit)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	// 返回多条数据
	for rows.Next() {
		tuf := TableUserFile{}
		err = rows.Scan(&tuf.FileName, &tuf.FileHash, &tuf.FileSize, &tuf.UploadAt, &tuf.LastUpdated)
		if err != nil {
			log.Println(err.Error())
			break
		}
		userFileArr = append(userFileArr, tuf)
	}
	res.Suc = true
	res.Data = userFileArr
	return
}

// 修改用户文件表的元信息
func UpdateUserFileInfoDB(username string, filename string, filehash string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("update tbl_user_file set file_name=? where user_name=? and file_sha1=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// 判断文件是否存在
func IsUserFileUpload(username string, filehash string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select 1 from tbl_user_file where user_name=? and file_sha1=? and status=1 limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil || !rows.Next() {
		res.Suc = true
		res.Data = map[string]bool{
			"exists": false,
		}
		return
	}
	res.Suc = true
	res.Data = map[string]bool{
		"exists": true,
	}
	return
}

// 用户文件名改之前，文件名有没有被使用
func QueryUserFileNameExist(username string, filename string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select 1 from tbl_user_file where user_name=? and file_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(username, filename)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil || !rows.Next() {
		res.Suc = true
		res.Data = map[string]bool{
			"exists": false,
		}
		return
	}
	res.Suc = true
	res.Data = map[string]bool{
		"exists": true,
	}
	return
}

// 查询单个用户文件元信息
func QueryUserFileDB(username string, filehash string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select file_name,file_sha1,file_size,upload_at,last_update from tbl_user_file where user_name=? and file_sha1=?")
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	ufile := TableUserFile{}
	for rows.Next() {
		err = rows.Scan(&ufile.FileName, &ufile.FileHash, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = ufile
	return
}

// 删除用户文件表
func DeleteUserFileDB(username string, filehash string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("update tbl_user_file set status=2 where user_name=? and file_sha1=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}
