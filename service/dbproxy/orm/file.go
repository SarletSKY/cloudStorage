package orm

import (
	"database/sql"
	_ "database/sql"
	"filestore-server-study/service/dbproxy/mysql"
	"log"
)

// 写入数据到mysql [文件上传到mysql完成]
func AddFileInfoTodb(fileSha1 string, fileName string, fileSize int64, fileAdd string) (res ExecResult) {
	// sql语句
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")

	if err != nil {
		log.Println("Failed to prepare statement add data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileSha1, fileName, fileSize, fileAdd)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	// 判断sha1是否已经存在
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Println("data already exist,", fileSha1)
		}
		res.Suc = true
		return
	}
	res.Suc = false
	return
}

// 更新数据到mysql
func UpdateFileInfoTodb(fileSha1 string, fileName string, fileSize int64, fileAdd string) (res ExecResult) {
	// sql语句
	stmt, err := mysql.DBConn().Prepare("update tbl_file set file_sha1=?,file_name=?,file_size=?,file_addr=? where file_sha1=?")

	if err != nil {
		log.Println("Failed to prepare statement update data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileSha1, fileName, fileSize, fileAdd, fileSha1)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	// 判断sha1是否已经存在
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Println("data already exist,", fileSha1)
		}
		res.Suc = true
		return
	}
	res.Suc = false
	return
}

// 从mysql获取数据到结构体
func GetFileInfoTodb(sha1 string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file " +
		"where file_sha1=? and status=1 limit 1")
	if err != nil {
		log.Println("Failed to prepare statement get data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var tf = TableFile{}
	// sql 获取查询 Scan()是将数据返回
	err = stmt.QueryRow(sha1).Scan(&tf.FileSha1, &tf.FileAdd, &tf.FileName, &tf.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应信息，返回参数错误为nil
			res.Suc = true
			res.Data = nil
			return
		} else {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = tf
	return
}

// 从mysql获取多条数据到结构体
func GetManyFileInfoTodb(count int) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file where status=1 limit ?")
	if err != nil {
		log.Println("Failed to prepare statement get many data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	tfArr := make([]TableFile, 0)
	rows, err := stmt.Query(count)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	// 返回多条数据 ??
	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	for i := 0; i < len(values) && rows.Next(); i++ {
		tf := TableFile{}
		err = rows.Scan(&tf.FileSha1, &tf.FileAdd, &tf.FileName, &tf.FileSize)
		if err != nil {
			log.Println(err.Error())
			break
		}
		tfArr = append(tfArr, tf)
	}
	res.Suc = true
	res.Data = tfArr
	log.Printf("find %d data\n", len(tfArr))
	return
}

// 删除[这里采用逻辑删除]
func DeleteFileMetaDB(sha1 string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set status=0 where file_sha1=?")
	if err != nil {
		log.Println("Failed to prepare statement delete data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(sha1)
	if err != nil {
		log.Println("sql exec failed,err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Println("data already exist,", sha1)
		}
		res.Suc = true
		return
	}
	res.Suc = false
	return
}

// 更新文件的路径
func UpdateFileLocationdb(filehash string, fileaddr string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set `file_addr`=? where `file_sha1`=? limit 1")
	if err != nil {
		log.Println("Failed to prepare statement get data ,err", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	// ??
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Printf("update data failed,:%s", filehash)
			res.Suc = false
			res.Msg = "无记录更新"
			return
		}
		res.Suc = true
		return
	}
	res.Suc = false
	res.Msg = err.Error()
	return
}
