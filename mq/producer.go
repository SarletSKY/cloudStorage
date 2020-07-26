package mq

import (
	"filestore-server-study/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// rabbitMQ的一个全局链接/全局通道/全局接收的err判断标志
var (
	conn        *amqp.Connection
	channel     *amqp.Channel
	notifyClose chan *amqp.Error //如果异常关闭，会接收通知
)

//异步处理重新链接
func init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel() {
		channel.NotifyClose(notifyClose)
	}

	//开启异步，断线重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel()
			}
		}
	}()
}

// 初始化rabbitMQ通道
func initChannel() bool {
	// 先判断是否已经初始化通道
	if channel != nil {
		return true
	}

	// 获取rabbitMQ链接
	var err error
	conn, err = amqp.Dial(config.RabbitURL)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// 生产者发布消息
func Publish(exchange string, routingKey string, msg []byte) bool {
	// 判断是否已经初始化通道
	if !initChannel() {
		return false
	}

	//发送消息
	if err := channel.Publish(
		exchange,
		routingKey,
		false, // 如果没有对应的queue，就自动弃掉这条消息
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	); err != nil {
		return false
	}
	return true
}
