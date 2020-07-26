package common

type StoreType int

const (
	_          StoreType = iota
	StoreLocal           // 单点本地
	StoreCeph            // ceph集群
	StoreOSS             // OSS
	StoreMix             // OSS加ceph
	StoreAll             // 所有类型都用
)
