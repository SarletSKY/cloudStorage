package main

import (
	"filestore-server-study/config"
	dbProto "filestore-server-study/service/dbproxy/proto"
	dbRPC "filestore-server-study/service/dbproxy/rpc"
	"github.com/micro/go-micro"
	"log"
	"time"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(config.RegistryConsul()),
	)
	service.Init()
	dbProto.RegisterDBProxyServiceHandler(service.Server(), new(dbRPC.DBProxy))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func main() {
	startRPCService()
}
