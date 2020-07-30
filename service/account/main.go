package main

import (
	"filestore-server-study/service/account/handler"
	proto "filestore-server-study/service/account/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"time"
)

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"172.17.0.5:8500",
		}
	})

	// 注册服务，创建一个service
	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),
		micro.Registry(reg),
		micro.RegisterInterval(time.Second*5),
	)
	// 初始化service
	service.Init()
	// 加入注册用户service
	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	// 运行
	err := service.Run()
	if err != nil {
		log.Println(err)
	}
}
