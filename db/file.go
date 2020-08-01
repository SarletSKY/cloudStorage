package db

import (
	"database/sql"
	"filestore-server-study/db/mysql"
	"fmt"
)

// 表的文件结构体
type TableFile struct {
	FileSha1 string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAdd  sql.NullString
}

// 写入数据到mysql [文件上传到mysql完成]
func AddFileInfoTodb(fileSha1 string, fileName string, fileSize int64, fileAdd string) bool {
	// sql语句
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")

	if err != nil {
		fmt.Println("Failed to prepare statement add data ,err", err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileSha1, fileName, fileSize, fileAdd)
	if err != nil {
		fmt.Println("sql exec failed,err: ", err.Error())
		return false
	}
	// 判断sha1是否已经存在
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("data already exist,", fileSha1)
		}
		return true
	}
	return false
}

// 更新数据到mysql
func UpdateFileInfoTodb(fileSha1 string, fileName string, fileSize int64, fileAdd string) bool {
	// sql语句
	stmt, err := mysql.DBConn().Prepare("update tbl_file set file_sha1=?,file_name=?,file_size=?,file_addr=? where file_sha1=?")

	if err != nil {
		fmt.Println("Failed to prepare statement update data ,err", err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileSha1, fileName, fileSize, fileAdd, fileSha1)
	if err != nil {
		fmt.Println("sql exec failed,err: ", err.Error())
		return false
	}
	// 判断sha1是否已经存在
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("data already exist,", fileSha1)
		}
		return true
	}
	return false
}

// 从mysql获取数据到结构体
func GetFileInfoTodb(sha1 string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file " +
		"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement get data ,err", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var tf = TableFile{}
	// sql 获取查询 Scan()是将数据返回
	err = stmt.QueryRow(sha1).Scan(&tf.FileSha1, &tf.FileAdd, &tf.FileName, &tf.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tf, nil
}

// 从mysql获取多条数据到结构体
func GetManyFileInfoTodb(count int) ([]*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file where status=1 limit ?")
	if err != nil {
		fmt.Println("Failed to prepare statement get many data ,err", err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfArr := make([]*TableFile, 0)
	rows, err := stmt.Query(count)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// 返回多条数据
	for rows.Next() {
		tf := TableFile{}
		err = rows.Scan(&tf.FileSha1, &tf.FileAdd, &tf.FileName, &tf.FileSize)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		tfArr = append(tfArr, &tf)
	}
	fmt.Printf("find %d data\n", len(tfArr))
	return tfArr, nil
}

// 删除[这里采用逻辑删除]
func DeleteFileMetaDB(sha1 string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set status=0 where file_sha1=?")
	if err != nil {
		fmt.Println("Failed to prepare statement delete data ,err", err.Error())
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(sha1)
	if err != nil {
		fmt.Println("sql exec failed,err: ", err.Error())
		return false
	}
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("data already exist,", sha1)
		}
		return true
	}
	return false
}

// 更新文件的路径
func UpdateFileLocationdb(filehash string, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set `file_addr`=? where `file_sha1`=? limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement get data ,err", err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("data already exist,", filehash)
		}
		return true
	}
	return false
}
