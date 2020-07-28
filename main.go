package main

import (
	"filestore-server-study/route"
	"fmt"
)

func main() {
	router := route.Router()

	err := router.Run()
	if err != nil {
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}
