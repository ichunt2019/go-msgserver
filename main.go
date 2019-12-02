package main

import (
	"fmt"
	"github.com/pkg/errors"
	"go-msgserver/utils/rabbitmq"
	"time"
)

type TestPro struct {
	msgContent   string
}

// 实现生产者
func (t *TestPro) MsgContent() string {
	return t.msgContent
}

//// 实现消费者
func (t *TestPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	return errors.New("消息处理失败")
	//return nil
}

func main() {
	msg := fmt.Sprintf("这是测试任务")
	t := &TestPro{
		msg,
	}
	queueExchange := &rabbitmq.QueueExchange{
		"b_test_rabbit",
		"b_test_rabbit",
		"b_test_rabbit_mq",
		"direct",
	}



	//生产消息
	//mq := rabbitmq.New(queueExchange)
	//mq.RegisterProducer(t)
	////for i :=0;i<=100 ;i++{
	////	mq.RegisterProducer(t)
	////	mq.RegisterProducer(t)
	////}
	//
	//mq.Start()
	//return



	//消费消息
	for{
		mq := rabbitmq.New(queueExchange)
		mq.RegisterReceiver(t)
		mq.Start()
		time.Sleep(time.Second*2)

	}

}