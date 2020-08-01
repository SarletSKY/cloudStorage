package handler

import (
	"context"
	"encoding/json"
	"filestore-server-study/common"
	"filestore-server-study/db"
	userProto "filestore-server-study/service/account/proto"
)

// 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *userProto.ReqUserFile, resp *userProto.RespUserFile) error {

	limit := int(req.Limit)
	username := req.Username

	userFile, err := db.QueryMantUserFileDB(username, limit)

	if err != nil {
		resp.Code = common.StatusServerError
		resp.Message = "查询用户文件错误"
		return nil
	}
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
	exist := db.QueryUserFileNameExist(username, newFileName)
	if exist {
		resp.Code = common.FileAlreadExists
		resp.Message = "文件名已经存在，请重新输入"
		return nil
	}

	_ = db.UpdateUserFileInfoDB(username, newFileName, filehash)

	// TODO: 6. 将用户文件表中更改的那条数据重新获取出来，序列化返回
	userFile, err := db.QueryUserFileDB(username, filehash)
	if err != nil {
		resp.Code = common.StatusServerError
		resp.Message = "更新用户文件名失败"
		return err
	}

	//fileMetaBytes, err := json.Marshal(curFileMeta)
	fileMetaBytes, err := json.Marshal(userFile)
	if err != nil {
		resp.Code = common.StatusMarshalInvalid
		resp.Message = "序列化失败"
		return err
	}
	if err != nil {
		resp.Code = common.StatusServerError
		return nil
	}

	resp.FileData = fileMetaBytes
	return nil
}
