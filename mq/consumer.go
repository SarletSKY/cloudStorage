package mq

import (
	"log"
)

var exit chan bool

func StartConsume(qName string, cName string, callback func(msg []byte) bool) {
	// 根据mq初始化生成消费者
	msg, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	exit = make(chan bool)

	// 循环读取msg
	go func() {
		for data := range msg { // 消费者接收队列消息[msg是chan型],传body数据
			// 处理msg信息
			suc := callback(data.Body)
			if suc {
				// TODO: 将任务写入错误队列，待后续处理
			}
		}
	}()

	// 为了让消费者一直执行，加入通道
	<-exit
	// 关闭通道
	channel.Close()
}

// 停止消费者
func StopConsume() {
	exit <- true
}
