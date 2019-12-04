package main

import (
	"fmt"
	_ "fmt"
	"github.com/ichunt2019/go-msgserver/utils/rabbitmq"
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
	err := mq.Start()
	if(err != nil){
		fmt.Println("发送消息失败")
	}








}