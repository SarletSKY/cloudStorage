package main

import (
	"filestore-server-study/config"
	upCfg "filestore-server-study/service/upload/config"
	upProto "filestore-server-study/service/upload/proto"
	"filestore-server-study/service/upload/route"
	upRpc "filestore-server-study/service/upload/rpc"
	"fmt"
	"github.com/micro/go-micro"
	"time"
)

// 启动rpc服务 micro服务
func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.upload"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(config.RegistryConsul()),
	)
	// 初始化服务
	service.Init()

	// 加入服务
	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))

	//日动服务
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// 启动api服务
func startAPIService() {
	router := route.Router()
	if err := router.Run(upCfg.UploadServiceHost); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// api服务
	go startAPIService()

	// rpc服务
	startRPCService()
}
