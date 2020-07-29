package rpc

import (
	"context"
	"filestore-server-study/service/upload/config"
	upProto "filestore-server-study/service/upload/proto"
)

type Upload struct{}

//获取上传入口
func (u *Upload) UploadEntry(ctx context.Context, req *upProto.ReqEntry, resp *upProto.RespEntry) error {
	resp.Entry = config.UploadEntry
	return nil
}
