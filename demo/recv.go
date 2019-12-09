package main

import (
	"fmt"
	"github.com/ichunt2019/go-msgserver/utils/rabbitmq"
	"time"
)

type RecvPro struct {

}

//// 实现消费者 消费消息失败 自动进入延时尝试  尝试3次之后入库db
func (t *RecvPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	//return errors.New("顶顶顶顶")
	return nil
}

//消息已经消费3次 失败了 请进行处理
func (t *RecvPro) FailAction(dataByte []byte) error {
	fmt.Println(string(dataByte))
	fmt.Println("任务处理失败了，我要进入db日志库了")
	fmt.Println("任务处理失败了，发送钉钉消息通知主人")
	return nil
}



func main() {

	//消费者实现 下面接口即可
	//type Receiver interface {
	//	Consumer([]byte)    error
	//}
	t := &RecvPro{}


	queueExchange := &rabbitmq.QueueExchange{
		"marking_export",
		"marking_export",
		"marking_export_exchange",
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