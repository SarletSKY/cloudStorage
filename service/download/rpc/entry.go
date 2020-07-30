package rpc

import (
	"context"
	"filestore-server-study/service/download/config"
	dlProto "filestore-server-study/service/download/proto"
)

type Download struct{}

func (d *Download) DownloadEntry(ctx context.Context, req *dlProto.ReqEntry, resp *dlProto.RespEntry) error {
	resp.Entry = config.DownloadEntry
	return nil
}
