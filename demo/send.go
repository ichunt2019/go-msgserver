package main

import (
	"fmt"
	"go-msgserver/utils/rabbitmq"
)

type SendPro struct {
	msgContent   string
}

// 实现生产者
func (t *SendPro) MsgContent() string {
	return t.msgContent
}

func main() {
	msg := fmt.Sprintf("这是测试任务")
	t := &SendPro{
		msg,
	}



	queueExchange := &rabbitmq.QueueExchange{
		"b_test_rabbit",
		"b_test_rabbit",
		"b_test_rabbit_mq",
		"direct",
	}
	mq := rabbitmq.New(queueExchange)
	mq.RegisterProducer(t)
	mq.Start()






}