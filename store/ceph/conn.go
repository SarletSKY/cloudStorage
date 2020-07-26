package ceph

import (
	cfg "filestore-server-study/config"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

// 设置全局链接
var cephConn *s3.S3

// 获取ceph链接
func GetCephConn() *s3.S3 {
	// 1. 没有链接就返回
	if cephConn != nil {
		return cephConn
	}

	// 2. 初始化ceph的信息
	auth := aws.Auth{
		AccessKey: cfg.CephAccessKey,
		SecretKey: cfg.CephSecretKey,
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          cfg.CephGWEndpoint,
		S3Endpoint:           cfg.CephGWEndpoint,
		S3LowercaseBucket:    false,
		S3LocationConstraint: false,
		S3BucketEndpoint:     "",
		Sign:                 aws.SignV2,
	}
	// 3. 创建s3类型的临界
	return s3.New(auth, curRegion)
}

// 获取指定Bucket对象
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConn()
	return conn.Bucket(bucket)
}

// PUtObject: 上传文件到ceph
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}
