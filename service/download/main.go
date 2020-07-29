package download

import (
	"filestore-server-study/config"
	dlCfg "filestore-server-study/service/download/config"
	dlProto "filestore-server-study/service/download/proto"
	"filestore-server-study/service/download/route"
	"filestore-server-study/service/download/rpc"
	"fmt"
	"github.com/micro/go-micro"
	"time"
)

// rpc服务
func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(config.RegistryConsul()),
	)
	service.Init()
	dlProto.RegisterDownloadServiceHandler(service.Server(), new(rpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}

}

// api服务
func startAPIService() {
	router := route.Router()
	router.Run(dlCfg.DownloadServiceHost)
}

func main() {
	go startAPIService()
	startRPCService()
}
