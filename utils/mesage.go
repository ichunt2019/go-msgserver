package utils

import "go-msgserver/utils/rabbitmq"


type IchuntMessage interface {
	Start()
	//mqClose()
	RegisterProducer(rabbitmq.Producer)
	RegisterReceiver(rabbitmq.Receiver)
	//listenProducer(interface{})
}
