package handler

import (
	"context"
	"encoding/json"
	"filestore-server-study/common"
	userProto "filestore-server-study/service/account/proto"
	dbCli "filestore-server-study/service/dbproxy/client"
)

// 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *userProto.ReqUserFile, resp *userProto.RespUserFile) error {

	limit := int(req.Limit)
	username := req.Username

	dbResp, err := dbCli.QueryMantUserFileDB(username, limit)
	if err != nil || !dbResp.Suc {
		resp.Code = common.StatusServerError
		return err
	}

	// 将数据转换成user_file表数据
	userFile := dbCli.ToTableUserFiles(dbResp.Data)

	// 序列化数据
	fileMetaBytes, err := json.Marshal(userFile)
	if err != nil {
		resp.Code = common.StatusMarshalInvalid
		resp.Message = "序列化失败"
		return nil
	}
	resp.FileData = fileMetaBytes
	return nil
}

// 用户文件重命名
func (u *User) UserFileRename(ctx context.Context, req *userProto.ReqUserFileRename, resp *userProto.RespUserFileRename) error {

	// 通过sha1获取文件的元信息 op是指客户端需要操作的类型的标志
	filehash := req.Filehash
	newFileName := req.NewFileName
	username := req.Username

	// TODO: 重命名之前查数据库有不有该名字，不能重复
	if exists, _ := dbCli.QueryUserFileNameExist(username, newFileName); exists {
		resp.Code = common.FileAlreadExists
		resp.Message = "文件名已经存在，请重新输入"
		return nil
	}

	dbResp, err := dbCli.UpdateUserFileInfoDB(username, newFileName, filehash)
	if err != nil || !dbResp.Suc {
		resp.Code = common.StatusServerError
		return err
	}

	userFile := dbCli.ToTableUserFiles(dbResp.Data)
	//fileMetaBytes, err := json.Marshal(curFileMeta)
	fileMetaBytes, err := json.Marshal(userFile)
	if err != nil {
		resp.Code = common.StatusMarshalInvalid
		resp.Message = "序列化失败"
		return err
	}

	resp.FileData = fileMetaBytes
	return nil
}
