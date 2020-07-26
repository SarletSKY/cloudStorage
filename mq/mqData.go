package mq

import "filestore-server-study/common"

// 将要写到rabbitMQ的数据结构体
type TransferData struct {
	FileHash      string
	CurLocation   string           // 本地路径
	DestLocation  string           // oss路径
	DestStoreType common.StoreType // 当前使用的存储类型，本地/ceph/oss
}
