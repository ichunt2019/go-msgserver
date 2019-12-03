package utils

import "go-msgserver/utils/rabbitmq"


type IchuntMessage interface {
	Start()
	RegisterProducer(rabbitmq.Producer)
	RegisterReceiver(rabbitmq.Receiver)
	//listenProducer(interface{})
}
