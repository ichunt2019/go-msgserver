package main

import (
	"fmt"
	_ "fmt"
	"github.com/ichunt2019/go-msgserver/utils/rabbitmq"
)

func main() {

	queueExchange := &rabbitmq.QueueExchange{
		"fengkong_dong_count",
		"fengkong_dong_count",
		"fengkong_exchange",
		"direct",
		"amqp://guest:guest@192.168.2.232:5672/",
	}
	mq := rabbitmq.New(queueExchange)
	for i := 0;i<1;i++{
		mq.RegisterProducer("{\"com_credits_id\":\"2\",\"erp_company_code\":\"LX001\"}")
	}
	err := mq.Start()
	if(err != nil){
		fmt.Println("发送消息失败")
	}








}