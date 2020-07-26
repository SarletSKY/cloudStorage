package config

// rabbitMQ配置
const (
	// 是否异步处理
	AsyncTransferEnable = true // true：异步
	// RabbitURL: 登录入口
	RabbitURL = "amqp://guest:guest@127.0.0.1:5672/"
	// 用于文件transfer的交换机
	TransExchangeName = "uploadserver.trans"
	// 用户bindKey与publicKey
	TransOSSRoutingKey = "oss"
	// oss转移的队列名
	TransOSSQueueName = "uploadserver.trans.oss"
	// oss转移出错的信息队列
	TransOSSErrQueueName = "uploadserver.trans.err"
)
