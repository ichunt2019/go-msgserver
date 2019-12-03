package main

import (
	"fmt"
	"github.com/pkg/errors"
	"go-msgserver/utils"
	_ "go-msgserver/utils/rabbitmq"
	"time"
	_ "time"
)

type TestPro struct {
	msgContent   string
}

// 实现生产者
func (t *TestPro) MsgContent() string {
	return t.msgContent
}

//// 实现消费者 消费消息失败 自动进入延时尝试  尝试3次之后入库db
func (t *TestPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	time.Sleep(time.Microsecond)
	return errors.New("消息处理失败")
	//return nil
}

func main() {
	msg := fmt.Sprintf("这是测试任务")
	t := &TestPro{
		msg,
	}


	utils.Init("rabbitmq","b_test_rabbit","b_test_rabbit","b_test_rabbit_mq")



	//生产消息
	//for i:=0;i<10;i++{
	//	utils.Ichunt2019MessageServer.RegisterProducer(t)
	//
	//}
	//utils.Ichunt2019MessageServer.Start()
	//return


	//消费消息
	for{
		utils.Ichunt2019MessageServer.RegisterReceiver(t)
		utils.Ichunt2019MessageServer.Start()
		time.Sleep(time.Second*1)

	}

}