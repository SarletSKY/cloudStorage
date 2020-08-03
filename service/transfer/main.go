package main

import (
	"bufio"
	"encoding/json"
	"filestore-server-study/config"
	"filestore-server-study/mq"
	dbCli "filestore-server-study/service/dbproxy/client"
	"filestore-server-study/store/oss"
	"fmt"
	"github.com/micro/go-micro"
	"log"
	"os"
	"time"
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
	resp, err := dbCli.UpdateFileLocationdb(transData.FileHash, transData.DestLocation)

	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !resp.Suc {
		log.Println("更新数据异常，请检查：" + transData.FileHash)
		return false
	}
	return true
}

// TODO: 将rabbitMQ 启动函数修改成微服务

//异步rabbitMQ
func startTransferService() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件被禁用，使用同步转移文件...")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")
	log.Println("开始监听转移消息队列")
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}

// rpc服务
func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.transfer"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(config.RegistryConsul()),
	)
	service.Init()
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
func main() {
	go startTransferService()

	startRPCService()
}
