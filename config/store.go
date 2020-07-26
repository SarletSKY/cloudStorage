package config

import "filestore-server-study/common"

const (
	TempLocalRootDir  = "/data/fileserver_tmp/" // 本地临时储存路径
	MergeLocalRootDir = "/data/fileserver_merge/"
	BlockLocalRootDir = "/data/fileserver_block" // 分块路径
	CurrentStoreType  = common.StoreOSS          // 设置当前的存储类型
)
