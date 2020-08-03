package client

import "time"

const layout = "2006-01-02 15:04:05"

type ByUploadTime []FileMeta

/*
	具体查官方文档
	Sort对data进行排序。它调用一次 data.Len 来决定排序的长度 n，调用 data.Less 和 data.Swap 的开销为 O(n*log(n))。此排序为不稳定排序。
	实际就是调用底层，将三个函数进行重写
*/
func (b ByUploadTime) Len() int {
	return len(b)
}

func (b ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(layout, b[i].UpdateFileTime)
	jTime, _ := time.Parse(layout, b[j].UpdateFileTime)
	return iTime.Nanosecond() > jTime.Nanosecond()
}

func (b ByUploadTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
