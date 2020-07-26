package oss

import (
	"filestore-server-study/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

func Client() *oss.Client {
	// oss为空时
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(config.OSSEndpoint, config.OSSAccesskeyID, config.OSSAccessKeySecret)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return ossCli
}

// 保存bucket存储空间
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(config.OSSBucket)
		if err != nil {
			return nil
		}
		return bucket
	}
	return nil
}

// downLoadURL: 临时授权下载url
func DownloadURL(objName string) string {
	signedURL, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return signedURL
}

// 针对制定Bucket设置生命周期
func BucketLifecycleRule(bucketName string) {
	// 表示前缀为test的对象，最后修改时间30天后过期
	rule1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rulesList := []oss.LifecycleRule{rule1}
	Client().SetBucketLifecycle(bucketName, rulesList)
}
