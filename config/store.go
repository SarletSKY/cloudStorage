package config

import "filestore-server-study/common"

const (
	TempLocalRootDir  = "/data/fileserver_tmp/"   // 本地临时储存路径
	MergeLocalRootDir = "/data/fileserver_merge/" // 合并路径
	BlockLocalRootDir = "/data/fileserver_block/" // 分块路径
	CephRootDir       = "/ceph"                   // ceph 私有云的文件路径
	OSSRootDir        = "oss/"                    // oss 公有云的文件路径
	CurrentStoreType  = common.StoreOSS           // 设置当前的存储类型
)
