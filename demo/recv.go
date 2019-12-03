package main

import (
	"fmt"
	"go-msgserver/utils/rabbitmq"
	"time"
)

type RecvPro struct {
	msgContent   string
}



//// 实现消费者 消费消息失败 自动进入延时尝试  尝试3次之后入库db
func (t *RecvPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	return nil
}

func main() {
	msg := fmt.Sprintf("这是测试任务")
	t := &RecvPro{
		msg,
	}

	

	queueExchange := &rabbitmq.QueueExchange{
		"b_test_rabbit",
		"b_test_rabbit",
		"b_test_rabbit_mq",
		"direct",
	}
	for{

		mq := rabbitmq.New(queueExchange)
		mq.RegisterReceiver(t)
		mq.Start()
		time.Sleep(time.Second)
	}





}