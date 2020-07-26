package test

import (
	"filestore-server-study/store/ceph"
	"fmt"
	"gopkg.in/amz.v1/s3"
	"os"
)

func main() {

	// 获取bucket
	bucket := ceph.GetCephBucket("userFile")

	// 获取数据进行测试
	data, _ := bucket.Get("/ceph/8601223fcd56ed69a21fcf643f7d0b9eba4ab64f")
	tmpFile, err := os.Create("/tmp/testfile")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tmpFile.Write(data)
	return

	// 创建一个新的bucket
	err = bucket.PutBucket(s3.PublicRead)
	if err != nil {
		fmt.Println("create bucket failed", err)
	}

	// 获取key值
	result, err := bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("bucket list err:%s\n", err)
	} else {
		fmt.Printf("object keys:%v\n", result)
	}

	// 上传文件
	objectPath := "/testupload/a.txt"
	err = bucket.Put(objectPath, []byte("put data"), "octet-stream", s3.PublicRead)
	if err != nil {
		fmt.Printf("upload data,err:%v\n", err)
	} else {
		fmt.Printf("object upload success")
	}

	// 在次获取key
	result, err = bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("bucket list err:%s\n", err)
	} else {
		fmt.Printf("object keys:%v\n", result)
	}
}
