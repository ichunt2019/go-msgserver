package main

import (
	_ "fmt"
	"go-msgserver/utils/rabbitmq"
	"strconv"
)

func main() {

	queueExchange := &rabbitmq.QueueExchange{
		"b_test_rabbit",
		"b_test_rabbit",
		"b_test_rabbit_mq",
		"direct",
		"amqp://guest:guest@192.168.2.232:5672/",
	}
	mq := rabbitmq.New(queueExchange)
	for i := 0;i<10;i++{
		mq.RegisterProducer("这是测试任务"+strconv.Itoa(i))
	}
	mq.Start()








}