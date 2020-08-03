package client

import (
	"context"
	"encoding/json"
	"filestore-server-study/config"
	"filestore-server-study/service/dbproxy/orm"
	dbProto "filestore-server-study/service/dbproxy/proto"
	"github.com/micro/go-micro"
	"github.com/mitchellh/mapstructure"
	"log"
)

// 元信息结构体
type FileMeta struct {
	FileSha1       string // Sha1[类似Id的作用]
	FileName       string // 文件名称
	FileSize       int64  // 文件大小
	UpdateFileTime string // 创建/更新文件时间
	Location       string // 文件路径
}

// rpc服务
var (
	dbCli dbProto.DBProxyService
)

// 客户端
func init() {
	service := micro.NewService(
		micro.Registry(config.RegistryConsul()),
	)
	service.Init()

	dbCli = dbProto.NewDBProxyService("go.micro.service.dbproxy", service.Client())
}

// 调用表的结构体
func TableFileToFileMeta(tfile orm.TableFile) FileMeta {
	return FileMeta{
		FileName: tfile.FileName.String,
		FileSha1: tfile.FileSha1,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAdd.String,
	}
}

// execAction: 向dbproxy请求执行action
func execAction(funcName string, paramJson []byte) (*dbProto.RespExec, error) {
	return dbCli.ExecuteAction(context.TODO(), &dbProto.ReqExec{
		Action: []*dbProto.SingleAction{
			&dbProto.SingleAction{
				Name:   funcName,
				Params: paramJson,
			},
		},
	})
}

// parseBody: 转换rpc返回的结果
func parseBody(resp *dbProto.RespExec) *orm.ExecResult {
	if resp == nil || resp.Data == nil {
		return nil
	}
	respList := []orm.ExecResult{}
	_ = json.Unmarshal(resp.Data, &respList)
	if len(respList) > 0 {
		return &respList[0]
	}
	return nil
}

// 转换用户表
func ToTableUser(src interface{}) orm.TableUser {
	user := orm.TableUser{}
	mapstructure.Decode(src, &user)
	return user
}

// 转换文件表
func ToTableFile(src interface{}) orm.TableFile {
	file := orm.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

func ToTableFiles(src interface{}) []orm.TableFile {
	file := []orm.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

func ToTableUserFile(src interface{}) orm.TableUserFile {
	userFile := orm.TableUserFile{}
	mapstructure.Decode(src, &userFile)
	return userFile
}

func ToTableUserFiles(src interface{}) []orm.TableUserFile {
	userFiles := []orm.TableUserFile{}
	mapstructure.Decode(src, &userFiles)
	return userFiles
}

// 获取元信息
func GetFileInfoTodb(filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash})
	res, err := execAction("/file/GetFileInfoTodb", uInfo)
	return parseBody(res), err
}

// 从mysql批量获取文件元信息
func GetManyFileInfoTodb(count int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{count})
	res, err := execAction("/file/GetManyFileInfoTodb", uInfo)
	return parseBody(res), err
}

// 文件上传完成，保存meta  // 这里改成结构体
func AddFileInfoTodb(fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location})
	res, err := execAction("/file/AddFileInfoTodb", uInfo)
	return parseBody(res), err
}

// 更新文件数据到mysql
func UpdateFileInfoTodb(fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location})
	res, err := execAction("/file/UpdateFileInfoTodb", uInfo)
	return parseBody(res), err
}

func UpdateFileLocationdb(filehash string, location string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash, location})
	res, err := execAction("/file/UpdateFileLocationdb", uInfo)
	return parseBody(res), err
}

func DeleteFileMetaDB(filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash})
	res, err := execAction("/file/DeleteFileMetaDB", uInfo)
	return parseBody(res), err
}

func SignInUserdb(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := execAction("/user/SignInUserdb", uInfo)
	return parseBody(res), err
}

func SignUpUserdb(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := execAction("/user/SignUpUserdb", uInfo)
	return parseBody(res), err
}

func QueryUserInfodb(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/QueryUserInfodb", uInfo)
	return parseBody(res), err
}

func RegisterTokendb(username string, token string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, token})
	res, err := execAction("/user/RegisterTokendb", uInfo)
	return parseBody(res), err
}

func QueryUserFileDB(username string, filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash})
	res, err := execAction("/userfile/QueryUserFileDB", uInfo)
	return parseBody(res), err
}

func QueryMantUserFileDB(username string, count int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, count})
	res, err := execAction("/userfile/QueryMantUserFileDB", uInfo)
	return parseBody(res), err
}

func OnUserFileUploadFinshedDB(username string, fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, fmeta.FileName, fmeta.FileSha1, fmeta.FileSize, fmeta.UpdateFileTime})
	res, err := execAction("/userfile/OnUserFileUploadFinshedDB", uInfo)
	return parseBody(res), err
}

func UpdateUserFileInfoDB(username string, filename string, filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filename, filehash})
	res, err := execAction("/userfile/UpdateUserFileInfoDB", uInfo)
	return parseBody(res), err
}

func QueryUserFileNameExist(username, filename string) (bool, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filename})
	res, err := execAction("/userfile/QueryUserFileNameExist", uInfo)
	if err != nil {
		return false, nil
	}

	execRes := parseBody(res)
	if execRes == nil {
		return false, nil
	}

	var data map[string]bool
	err = mapstructure.Decode(execRes.Data, &data)
	if err != nil {
		return false, err
	}
	log.Printf("QueryUserFileNameExist: %s %+v\n", username, data)
	return data["exists"], nil
}

func GetUserToken(username string) (string, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/GetUserToken", uInfo)

	if err != nil {
		return "", err
	}

	execRes := parseBody(res)
	if execRes == nil {
		return "", nil
	}

	var data map[string]string
	err = mapstructure.Decode(execRes.Data, &data)
	if err != nil {
		return "", err
	}

	log.Printf("GetUserToken: %+v\n", data)
	return data["token"], nil
}

func IsUserFileUpload(username, filehash string) (bool, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash})
	res, err := execAction("/userfile/IsUserFileUpload", uInfo)
	if err != nil {
		return false, err
	}

	execRes := parseBody(res)
	if execRes == nil {
		return false, nil
	}

	var data map[string]bool
	err = mapstructure.Decode(execRes.Data, &data)
	if err != nil {
		return false, err
	}
	log.Printf("IsUserFileUploaded: %s %s %+v\n", username, filehash, data)
	return data["exists"], nil
}

func UserExist(username string) (bool, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/UserExist", uInfo)
	if err != nil {
		return false, err
	}

	execRes := parseBody(res)
	if execRes == nil {
		return false, nil
	}

	var data map[string]bool
	err = mapstructure.Decode(execRes.Data, &data)
	if err != nil {
		return false, err
	}
	log.Printf("User already exists: %s", username)
	return data["exists"], nil
}
