package meta

import (
	"filestore-server-study/db"
	"sort"
)

// 定义一个文件元信息的结构提

type FileMeta struct {
	FileSha1       string // Sha1[类似Id的作用]
	FileName       string // 文件名称
	FileSize       int64  // 文件大小
	UpdateFileTime string // 创建/更新文件时间
	Location       string // 文件路径
}

// 将元信息存储起来
var fileMetaMap map[string]FileMeta

// 对文件进行初始化
func init() {
	fileMetaMap = make(map[string]FileMeta)
}

// 新增/更新元信息到map中
func UploadFileMeta(fm FileMeta) {
	fileMetaMap[fm.FileSha1] = fm
}

// 获取元信息
func GetFileMeta(sha1 string) FileMeta {
	return fileMetaMap[sha1]
}

// 查询多个元信息,业务实现通过输入数量来返回查询个数
func GetLastFileMeta(count int) []FileMeta {
	// 定义一个临时变量
	var fileMetaMapTemp []FileMeta

	for _, v := range fileMetaMap {
		fileMetaMapTemp = append(fileMetaMapTemp, v)
	}
	// 根据count来返回临时的变量
	sort.Sort(ByUploadTime(fileMetaMapTemp))

	// 如果输入大于总数值，则全部返回
	if count > len(fileMetaMapTemp) {
		return fileMetaMapTemp
	}
	return fileMetaMapTemp[0:count]
}

// 删除元信息 这里应该枷锁
func DeleteFileMeta(sha1 string) {
	delete(fileMetaMap, sha1)
}

/**
DB操作
*/

// 新增/更新元信息到mysql [数据库]
func UploadFileMetaDB(fm FileMeta) bool {
	return db.AddFileInfoTodb(fm.FileSha1, fm.FileName, fm.FileSize, fm.Location)
}

// 更新元信息 [数据库]
func UpdateFileMetaDB(fm FileMeta) bool {
	return db.UpdateFileInfoTodb(fm.FileSha1, fm.FileName, fm.FileSize, fm.Location)
}

// 获取元信息[数据库]
func GetFileMetaDB(sha1 string) (FileMeta, error) {
	tableFile, err := db.GetFileInfoTodb(sha1)
	if err != nil {
		return FileMeta{}, err
	}
	fileMeta := FileMeta{
		FileSha1: tableFile.FileSha1,
		FileName: tableFile.FileName.String,
		FileSize: tableFile.FileSize.Int64,
		Location: tableFile.FileAdd.String,
	}
	return fileMeta, nil
}

// 查询多个元信息,业务实现通过输入数量来返回查询个数 [数据库]
func GetLastFileMetaDB(count int) ([]FileMeta, error) {
	tfArrs, err := db.GetManyFileInfoTodb(count)
	if err != nil {
		return make([]FileMeta, 0), err
	}
	fileMetaMaps := make([]FileMeta, len(tfArrs))
	for i := 0; i < len(tfArrs); i++ {
		fileMetaMaps[i] = FileMeta{
			FileSha1: tfArrs[i].FileSha1,
			FileName: tfArrs[i].FileName.String,
			FileSize: tfArrs[i].FileSize.Int64,
			Location: tfArrs[i].FileAdd.String,
		}
	}
	return fileMetaMaps, nil
}

// 删除元信息 [数据路]
func DeleteFileMetaDB(sha1 string) bool {
	return db.DeleteFileMetaDB(sha1)
}
