package main

import (
	"bufio"
	"encoding/json"
	"filestore-server-study/config"
	"filestore-server-study/db"
	"filestore-server-study/mq"
	"filestore-server-study/store/oss"
	"log"
	"os"
)

// 实际消费者调用的文件转移函数
func ProcessTransfer(msg []byte) bool {
	// 解析获取queue的数据
	log.Println(string(msg))

	// 解析mq队列返回的消息
	transData := mq.TransferData{}
	// 将获取到的消息进行反序列化解析
	err := json.Unmarshal(msg, &transData)
	if err != nil {
		return false
	}

	localFile, err := os.Open(transData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// 传到oss bucket
	err = oss.Bucket().PutObject(
		transData.DestLocation,
		bufio.NewReader(localFile),
	)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// 更改用户表的数据
	_ = db.UpdateFileLocationdb(transData.FileHash, transData.DestLocation)

	return true
}

func main() {
	log.Println("开始监听转移消息队列")
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}
