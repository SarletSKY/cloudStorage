package main

import (
	"filestore-server-study/service/apigw/route"
	"fmt"
)

func main() {
	router := route.Router()
	fmt.Println("微服务micro框架启动：开始监听...")
	router.Run(":8080")
}
