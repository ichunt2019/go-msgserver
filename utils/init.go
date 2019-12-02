package utils

import (
	"fmt"
	"go-msgserver/utils/rabbitmq"
)

var (
	Ichunt2019MessageServer  IchuntMessage
)

func Init(mq string,args ...string) (err error){
	switch mq {
	case "rabbitmq":
		queueExchange := &rabbitmq.QueueExchange{
			args[0],
			args[1],
			args[2],
			"direct",
		}
		Ichunt2019MessageServer  = rabbitmq.New(queueExchange)
		return
	case "kafka":
		err = fmt.Errorf("not support kafka")
		return
	case "Nsq":
		err = fmt.Errorf("not support Nsq")
		return
	default:
		err = fmt.Errorf("not support")
		return
	}
}
