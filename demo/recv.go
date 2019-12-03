package main

import (
	"fmt"
	"go-msgserver/utils/rabbitmq"
	"time"
	_ "errors"
)

type RecvPro struct {

}

//// 实现消费者 消费消息失败 自动进入延时尝试  尝试3次之后入库db
func (t *RecvPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	//return errors.New("顶顶顶顶")
	return nil
}



func main() {

	//消费者实现 下面接口即可
	//type Receiver interface {
	//	Consumer([]byte)    error
	//}
	t := &RecvPro{}


	queueExchange := &rabbitmq.QueueExchange{
		"b_test_rabbit",
		"b_test_rabbit",
		"b_test_rabbit_mq",
		"direct",
		"amqp://guest:guest@192.168.2.232:5672/",
	}
	for{

		mq := rabbitmq.New(queueExchange)
		mq.RegisterReceiver(t)
		err :=mq.Start()
		if err != nil{

			fmt.Println(err)
		}
		time.Sleep(time.Second)
	}





}